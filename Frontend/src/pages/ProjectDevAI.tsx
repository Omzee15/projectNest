import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { GoogleGenerativeAI } from '@google/generative-ai';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { 
  ArrowLeft, Bot, User, Send, Loader2, RefreshCw, 
  MessageSquare, Zap, Plus, Trash2, Edit2, MoreHorizontal
} from 'lucide-react';
import { apiService } from '@/services/api';
import { toast } from '@/hooks/use-toast';
import { Project, ChatMessage, ChatConversation } from '@/types';

// Mock data for chat conversations - in a real app this would come from the backend
const mockConversations: ChatConversation[] = [];

const ProjectDevAI: React.FC = () => {
  const { projectUid } = useParams<{ projectUid: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  // State for conversations and chat
  const [conversations, setConversations] = useState<ChatConversation[]>(mockConversations);
  const [selectedConversation, setSelectedConversation] = useState<ChatConversation | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [genAI, setGenAI] = useState<GoogleGenerativeAI | null>(null);
  const [chatSession, setChatSession] = useState<any>(null);
  
  // UI State
  const [showCreateChatDialog, setShowCreateChatDialog] = useState(false);
  const [newChatName, setNewChatName] = useState('');
  const [editingConversation, setEditingConversation] = useState<ChatConversation | null>(null);
  const [editingChatName, setEditingChatName] = useState('');
  
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // System prompt from the DEV AI file
  const SYSTEM_PROMPT = `You are DevSprint-AI, an experienced senior developer, specializing in practical guidance.
Think and act like a seasoned engineer, not a passive assistant.

IMPORTANT CONTEXT: This is a learning project management tool. Users have existing projects with learning roadmaps, tasks, and goals. Your role is to:
- Help users understand and execute their learning tasks
- Provide practical guidance on completing project objectives
- Explain technical concepts related to their learning goals
- Suggest next steps for skill development
- Answer questions about methodologies and best practices

DO NOT assume users want to build new systems from scratch. They want help with their existing learning projects and tasks.

CODE GENERATION GUIDELINES:
- When users ask for scripts, code snippets, or implementation examples, ALWAYS provide the complete, working code
- Include all necessary imports, dependencies, and setup instructions
- Add helpful comments to explain the code
- Provide code in properly formatted code blocks with language syntax highlighting
- If the code is long, still provide it in full - don't just describe what the code should do
- For complex implementations, you may provide a brief explanation before or after the code

Communicate clearly, like a senior developer mentoring a peer — confident, practical, and direct.
Give detailed, complete answers when explaining concepts or providing code.
Give short and to the point answers only for simple questions, confirmations, or quick clarifications.`;

  // Fetch project data
  const { data: projectResponse, isLoading: projectLoading } = useQuery({
    queryKey: ['project', projectUid],
    queryFn: () => projectUid ? apiService.getProject(projectUid) : Promise.reject('No project ID'),
    enabled: !!projectUid,
  });

  const project = projectResponse?.data;

  // Initialize Gemini AI
  useEffect(() => {
    const initializeAI = async () => {
      try {
        const apiKey = import.meta.env.VITE_GEMINI_API_KEY;
        if (!apiKey) {
          toast({
            title: "API Key Missing",
            description: "Please add VITE_GEMINI_API_KEY to your environment variables",
            variant: "destructive",
          });
          return;
        }

        const ai = new GoogleGenerativeAI(apiKey);
        setGenAI(ai);

      } catch (error) {
        console.error('Failed to initialize AI:', error);
        toast({
          title: "AI Initialization Failed",
          description: "Failed to connect to Gemini AI",
          variant: "destructive",
        });
      }
    };

    if (project) {
      initializeAI();
    }
  }, [project]);

  // Initialize chat session when conversation is selected
  useEffect(() => {
    const initializeChatSession = async () => {
      if (!genAI || !selectedConversation || !project) return;

      try {
        const model = genAI.getGenerativeModel({ model: "gemini-2.0-flash" });

        // Create enhanced system prompt with project context
        const enhancedSystemPrompt = `${SYSTEM_PROMPT}

PROJECT CONTEXT:
- Project Name: ${project.name}
- Project Description: ${project.description || 'No description provided'}
- Project Status: ${project.status}
- Project Type: Learning roadmap and task management

This user has an existing learning project with organized lists and tasks. Help them understand their learning path, complete tasks effectively, and progress through their roadmap.`;

        // Convert existing messages to chat history format
        const history = selectedConversation.messages.length > 0 ? [
          {
            role: "user" as const,
            parts: [{ text: "Initialize with system prompt" }],
          },
          {
            role: "model" as const,
            parts: [{ text: `I understand. I'm DevSprint-AI, ready to help you with ${project.name}. I'll provide senior-level guidance on architecture, implementation, and best practices.` }],
          },
          ...selectedConversation.messages.map(msg => ({
            role: msg.type === 'user' ? 'user' as const : 'model' as const,
            parts: [{ text: msg.content }],
          }))
        ] : [
          {
            role: "user" as const,
            parts: [{ text: "Initialize with system prompt" }],
          },
          {
            role: "model" as const,
            parts: [{ text: `👋 **DevSprint-AI** here! I'm your learning mentor for **${project.name}**.

${project.description ? `I see you're working on: *${project.description}*` : ''}

I'm here to help you with:
- 📚 **Understanding your learning tasks** and how to approach them
- 🎯 **Breaking down complex topics** into manageable steps
- �️ **Practical guidance** on tools and techniques
- � **Next steps** for progressing through your roadmap
- � **Best practices** and real-world application tips

**💡 Pro Tips:**
- Use \`@tasks\` in your messages to include your complete learning roadmap
- Use \`@db\` in your messages to include your project's database schema
- Use \`@flow\` to include your project's flowchart/workflow

**Ready to dive into your learning journey?** Ask me about any task in your project, or say "get started" for guidance on where to begin!` }],
          },
        ];

        const chat = model.startChat({
          history,
          generationConfig: {
            temperature: 0.7,
            topK: 40,
            topP: 0.8,
            maxOutputTokens: 2048,
          },
        });

        setChatSession(chat);
        setMessages(selectedConversation.messages);

      } catch (error) {
        console.error('Failed to initialize chat session:', error);
        toast({
          title: "Chat Session Failed",
          description: "Failed to initialize chat session",
          variant: "destructive",
        });
      }
    };

    initializeChatSession();
  }, [genAI, selectedConversation, project]);

  // Auto-scroll to bottom when new messages are added
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Load conversations from localStorage
  useEffect(() => {
    if (projectUid) {
      const stored = localStorage.getItem(`devai_conversations_${projectUid}`);
      if (stored) {
        try {
          const parsedConversations = JSON.parse(stored).map((conv: any) => ({
            ...conv,
            created_at: new Date(conv.created_at),
            updated_at: new Date(conv.updated_at),
            messages: conv.messages.map((msg: any) => ({
              ...msg,
              timestamp: new Date(msg.timestamp)
            }))
          }));
          setConversations(parsedConversations);
        } catch (error) {
          console.error('Failed to load conversations:', error);
        }
      }
    }
  }, [projectUid]);

  // Save conversations to localStorage
  const saveConversations = (convs: ChatConversation[]) => {
    if (projectUid) {
      localStorage.setItem(`devai_conversations_${projectUid}`, JSON.stringify(convs));
    }
  };

  const handleCreateChat = () => {
    if (!newChatName.trim() || !projectUid) return;

    const newConversation: ChatConversation = {
      id: Date.now().toString(),
      project_uid: projectUid,
      title: newChatName.trim(),
      messages: [],
      created_at: new Date(),
      updated_at: new Date(),
      is_active: true,
    };

    const updatedConversations = [newConversation, ...conversations].slice(0, 10); // Keep only 10 conversations
    setConversations(updatedConversations);
    saveConversations(updatedConversations);
    setSelectedConversation(newConversation);
    setNewChatName('');
    setShowCreateChatDialog(false);

    toast({
      title: 'Chat created',
      description: 'Your new chat conversation has been created.',
    });
  };

  const handleDeleteConversation = (conversationId: string) => {
    const updatedConversations = conversations.filter(conv => conv.id !== conversationId);
    setConversations(updatedConversations);
    saveConversations(updatedConversations);
    
    if (selectedConversation?.id === conversationId) {
      setSelectedConversation(null);
      setMessages([]);
    }

    toast({
      title: 'Chat deleted',
      description: 'Your chat conversation has been deleted.',
    });
  };

  const handleEditConversation = (conversation: ChatConversation) => {
    setEditingConversation(conversation);
    setEditingChatName(conversation.title);
  };

  const handleUpdateConversation = () => {
    if (!editingConversation || !editingChatName.trim()) return;

    const updatedConversations = conversations.map(conv =>
      conv.id === editingConversation.id
        ? { ...conv, title: editingChatName.trim(), updated_at: new Date() }
        : conv
    );

    setConversations(updatedConversations);
    saveConversations(updatedConversations);
    
    if (selectedConversation?.id === editingConversation.id) {
      setSelectedConversation(prev => prev ? { ...prev, title: editingChatName.trim() } : null);
    }

    setEditingConversation(null);
    setEditingChatName('');

    toast({
      title: 'Chat updated',
      description: 'Your chat title has been updated.',
    });
  };

  // Process special commands like @db, @flow, and @tasks
  const processMessageContent = async (content: string): Promise<string> => {
    let processedContent = content;

    // Check for @db command
    if (content.includes('@db')) {
      if (project?.dbml_content) {
        const dbSchema = `\n\n**PROJECT DATABASE SCHEMA:**\n\`\`\`sql\n${project.dbml_content}\n\`\`\`\n`;
        processedContent = processedContent.replace(/@db/g, dbSchema);
      } else {
        processedContent = processedContent.replace(/@db/g, '\n\n**No database schema available for this project.**\n');
      }
    }

    // Check for @flow command  
    if (content.includes('@flow')) {
      if (project?.flowchart_content) {
        const flowchart = `\n\n**PROJECT FLOWCHART:**\n\`\`\`mermaid\n${project.flowchart_content}\n\`\`\`\n`;
        processedContent = processedContent.replace(/@flow/g, flowchart);
      } else {
        processedContent = processedContent.replace(/@flow/g, '\n\n**No flowchart available for this project.**\n');
      }
    }

    // Check for @tasks command
    if (content.includes('@tasks')) {
      if (project?.lists && project.lists.length > 0) {
        let tasksOverview = `\n\n**PROJECT LEARNING ROADMAP:**\n`;
        tasksOverview += `**Project:** ${project.name}\n\n`;
        
        project.lists.forEach((list: any, index: number) => {
          tasksOverview += `**${index + 1}. ${list.name}**\n`;
          if (list.tasks && list.tasks.length > 0) {
            list.tasks.forEach((task: any, taskIndex: number) => {
              tasksOverview += `   ${taskIndex + 1}. ${task.title} (${task.priority} priority)\n`;
              if (task.description) {
                tasksOverview += `      - ${task.description}\n`;
              }
            });
          } else {
            tasksOverview += `   - No tasks yet\n`;
          }
          tasksOverview += `\n`;
        });
        
        processedContent = processedContent.replace(/@tasks/g, tasksOverview);
      } else {
        processedContent = processedContent.replace(/@tasks/g, '\n\n**No tasks available for this project.**\n');
      }
    }

    return processedContent;
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputMessage(e.target.value);
  };

  // Check if message contains special commands
  const hasDbCommand = inputMessage.includes('@db');
  const hasFlowCommand = inputMessage.includes('@flow');
  const hasTasksCommand = inputMessage.includes('@tasks');

  const sendMessage = async () => {
    if (!inputMessage.trim() || !chatSession || isLoading || !selectedConversation) return;

    // Process special commands before sending
    const processedContent = await processMessageContent(inputMessage.trim());

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: inputMessage.trim(), // Store original user input
      timestamp: new Date(),
    };

    const loadingMessage: ChatMessage = {
      id: (Date.now() + 1).toString(),
      type: 'ai',
      content: '',
      timestamp: new Date(),
      isLoading: true,
    };

    const newMessages = [...messages, userMessage, loadingMessage];
    setMessages(newMessages);
    setInputMessage('');
    setIsLoading(true);

    try {
      // Send the processed content to the AI (with @db/@flow replaced)
      const result = await chatSession.sendMessage(processedContent);
      const response = await result.response;
      const text = response.text();

      const aiMessage: ChatMessage = {
        id: (Date.now() + 2).toString(),
        type: 'ai',
        content: text,
        timestamp: new Date(),
      };

      const finalMessages = [...messages, userMessage, aiMessage];
      setMessages(finalMessages);

      // Update the conversation
      const updatedConversation = {
        ...selectedConversation,
        messages: finalMessages,
        updated_at: new Date(),
      };

      const updatedConversations = conversations.map(conv =>
        conv.id === selectedConversation.id ? updatedConversation : conv
      );

      setConversations(updatedConversations);
      saveConversations(updatedConversations);
      setSelectedConversation(updatedConversation);

    } catch (error) {
      console.error('Failed to send message:', error);
      
      const errorMessage: ChatMessage = {
        id: (Date.now() + 2).toString(),
        type: 'ai',
        content: '❌ Sorry, I encountered an error processing your message. Please try again.',
        timestamp: new Date(),
      };

      const finalMessages = [...messages, userMessage, errorMessage];
      setMessages(finalMessages);

      toast({
        title: "Message Failed",
        description: "Failed to get AI response. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  if (projectLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p>Project not found</p>
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b px-4 py-3 flex items-center gap-4">
        <Button 
          variant="outline" 
          size="sm"
          onClick={() => navigate(`/project/${projectUid}`)}
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Project
        </Button>
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded bg-gradient-to-br from-blue-500 to-purple-600">
            <Bot className="w-4 h-4 text-white" />
          </div>
          <div>
            <h1 className="font-semibold">{project.name} - DevSprint-AI</h1>
            <p className="text-sm text-gray-600">Senior Development Assistant</p>
          </div>
        </div>
        <div className="ml-auto">
          <Badge variant="secondary" className="flex items-center gap-1">
            <Zap className="h-3 w-3" />
            Gemini 2.0 Flash
          </Badge>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <div className="w-80 bg-white border-r flex flex-col">
          {/* Sidebar Header */}
          <div className="p-4 border-b">
            <div className="flex items-center justify-between mb-3">
              <h2 className="font-medium">AI Conversations</h2>
              <div className="flex gap-1">
                <Dialog open={showCreateChatDialog} onOpenChange={setShowCreateChatDialog}>
                  <DialogTrigger asChild>
                    <Button variant="ghost" size="sm">
                      <Plus className="w-4 h-4" />
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Create New Chat</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Input
                        value={newChatName}
                        onChange={(e) => setNewChatName(e.target.value)}
                        placeholder="Chat name (e.g., 'API Design Discussion')"
                        onKeyDown={(e) => e.key === 'Enter' && handleCreateChat()}
                      />
                      <Button onClick={handleCreateChat} disabled={!newChatName.trim()}>
                        Create Chat
                      </Button>
                    </div>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
            <p className="text-xs text-gray-500">
              Max 10 conversations per project
            </p>
          </div>

          {/* Edit Conversation Dialog */}
          <Dialog open={!!editingConversation} onOpenChange={(open) => {
            if (!open) {
              setEditingConversation(null);
              setEditingChatName('');
            }
          }}>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Chat Name</DialogTitle>
              </DialogHeader>
              <div className="space-y-4">
                <Input
                  value={editingChatName}
                  onChange={(e) => setEditingChatName(e.target.value)}
                  placeholder="Chat name"
                  onKeyDown={(e) => e.key === 'Enter' && handleUpdateConversation()}
                />
                <div className="flex gap-2">
                  <Button onClick={handleUpdateConversation} disabled={!editingChatName.trim()}>
                    Update Chat
                  </Button>
                  <Button variant="outline" onClick={() => {
                    setEditingConversation(null);
                    setEditingChatName('');
                  }}>
                    Cancel
                  </Button>
                </div>
              </div>
            </DialogContent>
          </Dialog>

          {/* Sidebar Content */}
          <div className="flex-1 overflow-y-auto p-2">
            {conversations.length === 0 ? (
              <div className="text-center py-8">
                <Bot className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                <p className="text-sm text-gray-500 mb-3">No conversations yet</p>
                <Button 
                  variant="outline" 
                  size="sm"
                  onClick={() => setShowCreateChatDialog(true)}
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Start First Chat
                </Button>
              </div>
            ) : (
              <div className="space-y-1">
                {conversations.map((conversation) => (
                  <div
                    key={conversation.id}
                    className={`group flex items-center gap-2 px-3 py-2 rounded-md hover:bg-blue-50 cursor-pointer text-sm transition-all border ${
                      selectedConversation?.id === conversation.id
                        ? 'bg-blue-50 border-blue-200 text-blue-900'
                        : 'border-transparent hover:border-blue-200'
                    }`}
                    onClick={() => setSelectedConversation(conversation)}
                  >
                    <MessageSquare className="w-4 h-4 text-blue-500 flex-shrink-0" />
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">{conversation.title}</p>
                      <p className="text-xs text-gray-500">
                        {conversation.messages.length} messages • {conversation.updated_at.toLocaleDateString()}
                      </p>
                    </div>
                    <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleEditConversation(conversation);
                        }}
                      >
                        <Edit2 className="w-3 h-3" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 text-red-500 hover:text-red-600"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteConversation(conversation.id);
                        }}
                      >
                        <Trash2 className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Main Chat Area */}
        <div className="flex-1 flex flex-col">
          {selectedConversation ? (
            <>
              {/* Chat Header */}
              <div className="bg-white border-b px-6 py-4 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <MessageSquare className="w-5 h-5 text-blue-500" />
                  <div>
                    <h3 className="font-semibold">{selectedConversation.title}</h3>
                    <p className="text-sm text-gray-500">
                      {messages.length} messages
                    </p>
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setMessages([]);
                    const updatedConversation = {
                      ...selectedConversation,
                      messages: [],
                      updated_at: new Date(),
                    };
                    const updatedConversations = conversations.map(conv =>
                      conv.id === selectedConversation.id ? updatedConversation : conv
                    );
                    setConversations(updatedConversations);
                    saveConversations(updatedConversations);
                    setSelectedConversation(updatedConversation);
                  }}
                  className="flex items-center gap-2"
                >
                  <RefreshCw className="h-4 w-4" />
                  Clear Chat
                </Button>
              </div>

              {/* Messages Area */}
              <div className="flex-1 overflow-hidden">
                <ScrollArea className="h-full">
                  <div className="p-6 space-y-4">
                    {messages.length === 0 && (
                      <div className="text-center py-8">
                        <Bot className="w-12 h-12 text-blue-300 mx-auto mb-3" />
                        <p className="text-gray-600 mb-2">Start a conversation with DevSprint-AI</p>
                        <p className="text-sm text-gray-500">Ask about architecture, code reviews, or implementation strategies</p>
                      </div>
                    )}
                    
                    {messages.map((message) => (
                      <div
                        key={message.id}
                        className={`flex gap-3 ${
                          message.type === 'user' ? 'justify-end' : 'justify-start'
                        }`}
                      >
                        {message.type === 'ai' && (
                          <div className="flex-shrink-0">
                            <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center">
                              <Bot className="h-4 w-4 text-white" />
                            </div>
                          </div>
                        )}
                        
                        <div
                          className={`max-w-[80%] rounded-lg px-4 py-3 ${
                            message.type === 'user'
                              ? 'bg-blue-500 text-white ml-auto'
                              : 'bg-gray-100 text-gray-900'
                          }`}
                        >
                          {message.isLoading ? (
                            <div className="flex items-center gap-2">
                              <Loader2 className="h-4 w-4 animate-spin" />
                              <span className="text-sm">DevSprint-AI is thinking...</span>
                            </div>
                          ) : (
                            <div className="whitespace-pre-wrap text-sm">
                              {message.content}
                            </div>
                          )}
                          <div className={`text-xs mt-2 ${
                            message.type === 'user' ? 'text-blue-100' : 'text-gray-500'
                          }`}>
                            {message.timestamp.toLocaleTimeString()}
                          </div>
                        </div>

                        {message.type === 'user' && (
                          <div className="flex-shrink-0">
                            <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center">
                              <User className="h-4 w-4 text-white" />
                            </div>
                          </div>
                        )}
                      </div>
                    ))}
                    <div ref={messagesEndRef} />
                  </div>
                </ScrollArea>
              </div>

              {/* Input Area */}
              <div className="border-t bg-white p-4">
                {/* Special Commands Indicator */}
                {(hasDbCommand || hasFlowCommand || hasTasksCommand) && (
                  <div className="mb-3 flex flex-wrap gap-2">
                    {hasDbCommand && (
                      <div className="inline-flex items-center gap-1 bg-blue-50 border border-blue-200 text-blue-700 px-2 py-1 rounded-md text-xs font-medium">
                        <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                        @db - Database Schema Included
                      </div>
                    )}
                    {hasFlowCommand && (
                      <div className="inline-flex items-center gap-1 bg-green-50 border border-green-200 text-green-700 px-2 py-1 rounded-md text-xs font-medium">
                        <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                        @flow - Flowchart Included
                      </div>
                    )}
                    {hasTasksCommand && (
                      <div className="inline-flex items-center gap-1 bg-purple-50 border border-purple-200 text-purple-700 px-2 py-1 rounded-md text-xs font-medium">
                        <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
                        @tasks - Learning Roadmap Included
                      </div>
                    )}
                  </div>
                )}
                
                <div className="flex gap-2">
                  <Textarea
                    placeholder="Ask DevSprint-AI about your learning tasks... Use @tasks for roadmap, @db for schema, @flow for flowchart"
                    value={inputMessage}
                    onChange={handleInputChange}
                    onKeyPress={handleKeyPress}
                    className="flex-1 min-h-[60px] resize-none"
                    disabled={isLoading}
                  />
                  <Button
                    onClick={sendMessage}
                    disabled={isLoading || !inputMessage.trim()}
                    size="icon"
                    className="h-[60px] w-[60px] shrink-0"
                  >
                    {isLoading ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Send className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                <div className="flex justify-between items-center mt-2">
                  <p className="text-xs text-gray-500">
                    Press Enter to send, Shift+Enter for new line
                  </p>
                  <p className="text-xs text-blue-600">
                    💡 Try: <code className="bg-gray-100 px-1 rounded">@tasks</code> for roadmap, <code className="bg-gray-100 px-1 rounded">@db</code> for schema, <code className="bg-gray-100 px-1 rounded">@flow</code> for flowchart
                  </p>
                </div>
              </div>
            </>
          ) : (
            /* Empty state */
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <Bot className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">Select a chat to continue</h3>
                <p className="text-gray-600 mb-4">Choose a conversation from the sidebar or create a new one</p>
                <Button onClick={() => setShowCreateChatDialog(true)}>
                  <Plus className="w-4 h-4 mr-2" />
                  Start New Chat
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProjectDevAI;