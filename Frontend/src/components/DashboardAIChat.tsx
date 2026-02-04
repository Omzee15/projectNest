import React, { useState, useEffect, useRef } from 'react';
import { GoogleGenerativeAI, SchemaType } from '@google/generative-ai';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { 
  Bot, User, Send, Loader2, RefreshCw, 
  MessageSquare, Zap, Plus, Trash2, Edit2, MoreHorizontal, X, CheckCircle, AlertCircle
} from 'lucide-react';
import { toast } from '@/hooks/use-toast';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { apiService } from '@/services/api';
import { AIProjectCreationRequest } from '@/types';

// Local types for the main dashboard chat (user-specific)
interface DashboardChatMessage {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
  isLoading?: boolean;
}

interface DashboardChatConversation {
  id: string;
  title: string;
  messages: DashboardChatMessage[];
  created_at: Date;
  updated_at: Date;
  is_active: boolean;
}

interface DashboardAIChatProps {
  isOpen: boolean;
  onClose: () => void;
}

export const DashboardAIChat: React.FC<DashboardAIChatProps> = ({ isOpen, onClose }) => {
  // State for conversations and chat
  const [conversations, setConversations] = useState<DashboardChatConversation[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<DashboardChatConversation | null>(null);
  const [messages, setMessages] = useState<DashboardChatMessage[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [genAI, setGenAI] = useState<GoogleGenerativeAI | null>(null);
  const [chatSession, setChatSession] = useState<any>(null);
  
  // UI State
  const [showCreateChatDialog, setShowCreateChatDialog] = useState(false);
  const [newChatName, setNewChatName] = useState('');
  const [editingConversation, setEditingConversation] = useState<DashboardChatConversation | null>(null);
  const [editingChatName, setEditingChatName] = useState('');
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);

  // System prompt for dashboard AI
  const SYSTEM_PROMPT = `You are ProjectNest-AI, a helpful assistant for the ProjectNest project management platform.

CONTEXT: This is the main dashboard chat where users can ask about:
- General project management advice and best practices
- Help with using ProjectNest features and functionality
- Guidance on organizing projects, tasks, and workflows
- General productivity and learning tips
- Questions about project methodologies and frameworks

**SPECIAL CAPABILITY: PROJECT CREATION**
You can actually create complete projects for users when they request it. When a user asks to create a project, follow this process:

1. **Detection**: Look for phrases like "create project", "new project", "make a project", "set up project", "build project", "help me plan", etc.

2. **Confirmation**: Before creating, ALWAYS ask: "Should I go ahead and create the project roadmap for [project idea]?"

3. **Function Calling**: Only after user confirms with "yes", "go ahead", "create it", "proceed", etc., use the create_project function with:
   - project_name: Clear, concise project name
   - project_description: Brief description of the project
   - lists: Array of task lists with organized tasks (e.g., Planning, Development, Testing, Deployment)

The system will automatically create the complete project structure in the database with all components.

You are NOT tied to any specific project, but rather help users with their overall project management journey and platform usage.

Communicate clearly and helpfully, like a knowledgeable project management advisor — professional, supportive, and practical.
Give clear, actionable advice that helps users make the most of their project management experience.
Keep responses concise but informative, expanding when detailed explanations are needed.`;

  // Initialize Gemini AI
  useEffect(() => {
    const initializeAI = async () => {
      try {
        console.log('🚀 [DashboardAIChat] Initializing Gemini AI...');
        const apiKey = import.meta.env.VITE_GEMINI_API_KEY;
        
        if (!apiKey) {
          console.error('❌ [DashboardAIChat] API Key missing');
          toast({
            title: "API Key Missing",
            description: "Please add VITE_GEMINI_API_KEY to your environment variables",
            variant: "destructive",
          });
          return;
        }

        console.log('🔑 [DashboardAIChat] API Key found, creating GoogleGenerativeAI instance...');
        const ai = new GoogleGenerativeAI(apiKey);
        setGenAI(ai);
        console.log('✅ [DashboardAIChat] Gemini AI initialized successfully');

      } catch (error) {
        console.error('❌ [DashboardAIChat] Failed to initialize AI:', error);
        toast({
          title: "AI Initialization Failed",
          description: "Failed to connect to Gemini AI",
          variant: "destructive",
        });
      }
    };

    if (isOpen) {
      console.log('🎯 [DashboardAIChat] Dialog opened, initializing AI...');
      initializeAI();
    }
  }, [isOpen]);

  // Initialize chat session when conversation is selected
  useEffect(() => {
    const initializeChatSession = async () => {
      if (!genAI || !selectedConversation) {
        console.log('⏸️ [DashboardAIChat] Skipping chat session init - missing genAI or selectedConversation');
        return;
      }

      try {
        console.log('🔧 [DashboardAIChat] Initializing chat session for conversation:', selectedConversation.id);
        
        console.log('🛠️ [DashboardAIChat] Creating function declarations...');

        const model = genAI.getGenerativeModel({ 
          model: "gemini-2.5-flash",
          tools: [{
            functionDeclarations: [{
              name: 'create_project',
              description: 'Creates a new project with tasks, lists, and additional components. Use this when the user asks to create, build, or set up a project.',
              parameters: {
                type: SchemaType.OBJECT,
                properties: {
                  project_name: {
                    type: SchemaType.STRING,
                    description: 'A clear, concise name for the project'
                  },
                  project_description: {
                    type: SchemaType.STRING,
                    description: 'A brief description of what the project is about'
                  },
                  lists: {
                    type: SchemaType.ARRAY,
                    description: 'Array of task lists/categories for organizing project tasks',
                    items: {
                      type: SchemaType.OBJECT,
                      properties: {
                        name: {
                          type: SchemaType.STRING,
                          description: 'Name of the task list (e.g., Planning, Development, Testing)'
                        },
                        tasks: {
                          type: SchemaType.ARRAY,
                          description: 'Tasks within this list',
                          items: {
                            type: SchemaType.OBJECT,
                            properties: {
                              title: {
                                type: SchemaType.STRING,
                                description: 'Task title'
                              },
                              description: {
                                type: SchemaType.STRING,
                                description: 'Task description'
                              }
                            },
                            required: ['title']
                          }
                        }
                      },
                      required: ['name', 'tasks']
                    }
                  }
                },
                required: ['project_name', 'project_description', 'lists']
              }
            }]
          }]
        });

        console.log('🤖 [DashboardAIChat] Model created with function declarations');

        // Convert existing messages to chat history format
        const history = selectedConversation.messages.length > 0 ? [
          {
            role: "user" as const,
            parts: [{ text: "Initialize with system prompt" }],
          },
          {
            role: "model" as const,
            parts: [{ text: `I understand. I'm ProjectNest-AI, ready to help you with project management and platform guidance.` }],
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
            parts: [{ text: `👋 **ProjectNest-AI** here! I'm your project management assistant.

I'm here to help you with:
- 📊 **Project management** best practices and strategies
- 🛠️ **Platform features** and how to use ProjectNest effectively
- 📋 **Task organization** and workflow optimization
- 🎯 **Productivity tips** for managing multiple projects
- 💡 **General guidance** on project methodologies

**💡 Tips:**
- Ask me about organizing your projects
- Get advice on task prioritization
- Learn about different project management methodologies
- Get help with platform features

**Ready to boost your project management game?** What would you like to know about managing your projects?` }],
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
        console.log('✅ [DashboardAIChat] Chat session initialized successfully');

      } catch (error) {
        console.error('❌ [DashboardAIChat] Failed to initialize chat session:', error);
        console.error('❌ [DashboardAIChat] Error details:', {
          name: error.name,
          message: error.message,
          stack: error.stack
        });
        toast({
          title: "Chat Session Failed",
          description: `Failed to initialize chat session: ${error.message}`,
          variant: "destructive",
        });
      }
    };

    initializeChatSession();
  }, [genAI, selectedConversation]);

  // Auto-scroll to bottom when new messages are added
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Load conversations from localStorage
  useEffect(() => {
    if (isOpen) {
      const stored = localStorage.getItem('dashboard_ai_conversations');
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
  }, [isOpen]);

  // Save conversations to localStorage
  const saveConversations = (convs: DashboardChatConversation[]) => {
    localStorage.setItem('dashboard_ai_conversations', JSON.stringify(convs));
  };

  const handleCreateChat = () => {
    if (!newChatName.trim()) return;

    const newConversation: DashboardChatConversation = {
      id: Date.now().toString(),
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

  const handleEditConversation = (conversation: DashboardChatConversation) => {
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

  // Function to handle tool calls (project creation)
  const handleToolCall = async (conversationText: string): Promise<string> => {
    try {
      console.log('🔨 [DashboardAIChat] handleToolCall started');
      console.log('📄 [DashboardAIChat] Conversation text length:', conversationText.length);
      console.log('📝 [DashboardAIChat] Full conversation text:', conversationText);

      // Show loading state
      setIsLoading(true);
      console.log('⏳ [DashboardAIChat] Loading state set to true');

      // Call the API to create the project using the entire conversation
      console.log('🌐 [DashboardAIChat] Calling API to create project...');
      const response = await apiService.createProjectFromAI(conversationText);
      console.log('📨 [DashboardAIChat] API response received:', response);
      
      if (response.success) {
        const project = response.data;
        console.log('✅ [DashboardAIChat] Project created successfully:', project);
        
        toast({
          title: "Project Created Successfully!",
          description: `Created "${project.name}" - check your Projects page to view it.`,
          variant: "default",
        });
        
        return `✅ **Project Created Successfully!**

📊 **Project Details:**
- **Name:** ${project.name}
- **Description:** ${project.description || 'Project created from our conversation'}
- **Project ID:** ${project.project_uid}

🎉 **What's Next?**
Your project has been created and is now available on your **Projects** page! The AI has automatically:
- ✅ Created the project structure
- ✅ Added organized task lists 
- ✅ Generated relevant tasks
- ✅ Set up any additional components (flowchart, schema) if mentioned

You can now navigate to your **Projects** page to view and manage this new project!

Would you like me to help you plan any additional aspects of this project?`;
      } else {
        throw new Error('Failed to create project');
      }
    } catch (error) {
      console.error('Error in tool call:', error);
      toast({
        title: "Error Creating Project",
        description: "There was an error creating the project. Please try again.",
        variant: "destructive",
      });
      return `❌ **Error Creating Project**

I apologize, but I encountered an error while creating the project. This could be due to:
- Network connectivity issues
- Server temporary unavailability
- Invalid project data format

Please try again in a moment, or let me know if you'd like to modify the project details.`;
    } finally {
      setIsLoading(false);
    }
  };

  const handleSendMessage = async () => {
    console.log('📤 [DashboardAIChat] handleSendMessage called');
    
    if (!inputMessage.trim() || !chatSession || !selectedConversation) {
      console.log('⚠️ [DashboardAIChat] Missing required data:', {
        hasMessage: !!inputMessage.trim(),
        hasChatSession: !!chatSession,
        hasSelectedConversation: !!selectedConversation
      });
      return;
    }

    console.log('💬 [DashboardAIChat] Sending message:', inputMessage.trim());

    const userMessage: DashboardChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: inputMessage.trim(),
      timestamp: new Date(),
    };

    const loadingMessage: DashboardChatMessage = {
      id: (Date.now() + 1).toString(),
      type: 'ai',
      content: '',
      timestamp: new Date(),
      isLoading: true,
    };

    // Update messages and conversation
    const newMessages = [...selectedConversation.messages, userMessage, loadingMessage];
    setMessages(newMessages);
    
    // Update the conversation in the list
    const updatedConversations = conversations.map(conv =>
      conv.id === selectedConversation.id
        ? { ...conv, messages: newMessages, updated_at: new Date() }
        : conv
    );
    setConversations(updatedConversations);
    setSelectedConversation(prev => prev ? { ...prev, messages: newMessages } : null);

    setInputMessage('');
    setIsLoading(true);

    console.log('🔄 [DashboardAIChat] Updated UI, sending to Gemini...');

    try {
      console.log('🤖 [DashboardAIChat] Calling chatSession.sendMessage...');
      const result = await chatSession.sendMessage(inputMessage.trim());
      console.log('📨 [DashboardAIChat] Received response from Gemini');
      
      const response = result.response;

      // Check for function calls in the response
      console.log('🔍 [DashboardAIChat] Checking for function calls...');
      const functionCalls = response.functionCalls?.();
      console.log('🔧 [DashboardAIChat] Function calls found:', functionCalls?.length || 0);
      
      if (functionCalls && functionCalls.length > 0) {
        console.log('📋 [DashboardAIChat] Function calls details:', functionCalls);
      }

      let aiResponse = response.text();
      console.log('💭 [DashboardAIChat] AI response text:', aiResponse);

      if (functionCalls && functionCalls.length > 0) {
        console.log('✨ [DashboardAIChat] Processing function calls...');
        
        // Show a temporary message that we're creating the project
        const tempAiMessage: DashboardChatMessage = {
          id: loadingMessage.id,
          type: 'ai',
          content: '🚀 **Creating your project...** \n\nI\'m setting up your project with all the components you requested. This will just take a moment!',
          timestamp: new Date(),
          isLoading: false,
        };

        // Update messages with temp message
        const tempMessages = [...selectedConversation.messages, userMessage, tempAiMessage];
        setMessages(tempMessages);

        // Process function calls
        for (const functionCall of functionCalls) {
          console.log('🔧 [DashboardAIChat] Processing function call:', functionCall.name);
          
          if (functionCall.name === 'create_project') {
            console.log('🎯 [DashboardAIChat] create_project function detected!');
            console.log('📊 [DashboardAIChat] Function arguments:', functionCall.args);
            
            // Prepare the conversation context for the backend AI
            const conversationContext = newMessages.map(msg => 
              `${msg.type === 'user' ? 'User' : 'AI'}: ${msg.content}`
            ).join('\n\n');

            console.log('📝 [DashboardAIChat] Conversation context prepared:', conversationContext);

            // Execute the tool call with the entire conversation
            console.log('🚀 [DashboardAIChat] Executing tool call...');
            const toolResult = await handleToolCall(conversationContext);
            aiResponse = toolResult;
            console.log('✅ [DashboardAIChat] Tool call completed');
            break;
          }
        }
      } else {
        console.log('💬 [DashboardAIChat] No function calls detected, using direct response');
      }

      const aiMessage: DashboardChatMessage = {
        id: loadingMessage.id,
        type: 'ai',
        content: aiResponse,
        timestamp: new Date(),
        isLoading: false,
      };

      // Replace loading message with actual response
      const finalMessages = [...selectedConversation.messages, userMessage, aiMessage];
      setMessages(finalMessages);

      const finalConversations = conversations.map(conv =>
        conv.id === selectedConversation.id
          ? { ...conv, messages: finalMessages, updated_at: new Date() }
          : conv
      );
      setConversations(finalConversations);
      saveConversations(finalConversations);
      setSelectedConversation(prev => prev ? { ...prev, messages: finalMessages } : null);

    } catch (error) {
      console.error('❌ [DashboardAIChat] Error sending message:', error);
      console.error('❌ [DashboardAIChat] Error details:', {
        name: error.name,
        message: error.message,
        stack: error.stack
      });
      
      // Remove loading message on error
      const errorMessages = [...selectedConversation.messages, userMessage];
      setMessages(errorMessages);

      const errorConversations = conversations.map(conv =>
        conv.id === selectedConversation.id
          ? { ...conv, messages: errorMessages }
          : conv
      );
      setConversations(errorConversations);
      setSelectedConversation(prev => prev ? { ...prev, messages: errorMessages } : null);

      toast({
        title: "Error",
        description: `Failed to send message: ${error.message}`,
        variant: "destructive",
      });
    } finally {
      console.log('🏁 [DashboardAIChat] handleSendMessage completed, setting loading to false');
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-2 sm:p-4">
      <Card className="ai-chat-dialog bg-card border w-full max-w-6xl h-[95vh] sm:h-[80vh] flex flex-col shadow-2xl">
        <CardHeader className="flex-shrink-0 bg-card border-b p-3 sm:p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Bot className="h-5 w-5 sm:h-6 sm:w-6 text-primary" />
              <CardTitle className="text-base sm:text-lg">ProjectNest AI</CardTitle>
              <Badge variant="secondary" className="ml-2 hidden sm:flex">
                <Zap className="h-3 w-3 mr-1" />
                Powered by Gemini
              </Badge>
            </div>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        </CardHeader>

        <CardContent className="flex-1 flex flex-col sm:flex-row gap-2 sm:gap-4 p-2 sm:p-4 min-h-0 bg-card overflow-hidden">
          {/* Conversations Sidebar */}
          <div className="w-full sm:w-80 flex flex-col gap-2 sm:gap-4 bg-card max-h-32 sm:max-h-full overflow-y-auto sm:overflow-visible">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-xs sm:text-sm text-muted-foreground uppercase tracking-wide">
                Conversations
              </h3>
              <Dialog open={showCreateChatDialog} onOpenChange={setShowCreateChatDialog}>
                <DialogTrigger asChild>
                  <Button size="sm" variant="outline" className="text-xs">
                    <Plus className="h-3 w-3 sm:h-4 sm:w-4 sm:mr-1" />
                    <span className="hidden sm:inline">New</span>
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Create New Chat</DialogTitle>
                  </DialogHeader>
                  <div className="space-y-4">
                    <Input
                      placeholder="Enter chat name..."
                      value={newChatName}
                      onChange={(e) => setNewChatName(e.target.value)}
                      onKeyPress={(e) => e.key === 'Enter' && handleCreateChat()}
                    />
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" onClick={() => setShowCreateChatDialog(false)}>
                        Cancel
                      </Button>
                      <Button onClick={handleCreateChat} disabled={!newChatName.trim()}>
                        Create Chat
                      </Button>
                    </div>
                  </div>
                </DialogContent>
              </Dialog>
            </div>

            <ScrollArea className="flex-1">
              <div className="space-y-2">
                {conversations.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    <MessageSquare className="h-8 w-8 mx-auto mb-2 opacity-50" />
                    <p className="text-sm">No conversations yet</p>
                    <p className="text-xs">Create your first chat to get started</p>
                  </div>
                ) : (
                  conversations.map((conversation) => (
                    <div
                      key={conversation.id}
                      className={`group relative p-3 rounded-lg cursor-pointer transition-colors ${
                        selectedConversation?.id === conversation.id
                          ? 'bg-primary/10 border border-primary/20'
                          : 'hover:bg-muted/50'
                      }`}
                      onClick={() => setSelectedConversation(conversation)}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1 min-w-0">
                          <h4 className="font-medium text-sm truncate">
                            {conversation.title}
                          </h4>
                          <p className="text-xs text-muted-foreground mt-1">
                            {conversation.messages.length} messages
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {conversation.updated_at.toLocaleDateString()}
                          </p>
                        </div>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="opacity-0 group-hover:opacity-100 h-6 w-6 p-0"
                              onClick={(e) => e.stopPropagation()}
                            >
                              <MoreHorizontal className="h-3 w-3" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation();
                                handleEditConversation(conversation);
                              }}
                            >
                              <Edit2 className="h-4 w-4 mr-2" />
                              Rename
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation();
                                handleDeleteConversation(conversation.id);
                              }}
                              className="text-destructive focus:text-destructive"
                            >
                              <Trash2 className="h-4 w-4 mr-2" />
                              Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </ScrollArea>
          </div>

          <Separator orientation="vertical" className="hidden sm:block" />
          <Separator orientation="horizontal" className="sm:hidden" />

          {/* Chat Area */}
          <div className="flex-1 flex flex-col min-w-0 bg-card">
            {selectedConversation ? (
              <>
                {/* Chat Header */}
                <div className="flex items-center justify-between pb-4">
                  <div>
                    <h3 className="font-semibold">{selectedConversation.title}</h3>
                    <p className="text-sm text-muted-foreground">
                      {selectedConversation.messages.length} messages
                    </p>
                  </div>
                </div>

                {/* Messages */}
                <ScrollArea className="flex-1 mb-4">
                  <div className="space-y-3 sm:space-y-4 p-1">
                    {messages.map((message) => (
                      <div
                        key={message.id}
                        className={`flex gap-2 sm:gap-3 ${
                          message.type === 'user' ? 'justify-end' : 'justify-start'
                        }`}
                      >
                        {message.type === 'ai' && (
                          <div className="flex-shrink-0">
                            <div className="w-6 h-6 sm:w-8 sm:h-8 rounded-full bg-primary/10 flex items-center justify-center">
                              <Bot className="h-3 w-3 sm:h-4 sm:w-4 text-primary" />
                            </div>
                          </div>
                        )}
                        <div
                          className={`max-w-[85%] sm:max-w-[70%] rounded-lg px-3 py-2 sm:px-4 ${
                            message.type === 'user'
                              ? 'bg-primary text-primary-foreground'
                              : 'ai-message-bg'
                          }`}
                        >
                          {message.isLoading ? (
                            <div className="flex items-center gap-2">
                              <Loader2 className="h-4 w-4 animate-spin" />
                              <span>Thinking...</span>
                            </div>
                          ) : (
                            <div className="whitespace-pre-wrap text-sm">
                              {message.content}
                            </div>
                          )}
                          <div
                            className={`text-xs mt-1 ${
                              message.type === 'user' ? 'text-primary-foreground/70' : 'text-muted-foreground'
                            }`}
                          >
                            {message.timestamp.toLocaleTimeString()}
                          </div>
                        </div>
                        {message.type === 'user' && (
                          <div className="flex-shrink-0">
                            <div className="w-6 h-6 sm:w-8 sm:h-8 rounded-full bg-primary flex items-center justify-center">
                              <User className="h-3 w-3 sm:h-4 sm:w-4 text-primary-foreground" />
                            </div>
                          </div>
                        )}
                      </div>
                    ))}
                    <div ref={messagesEndRef} />
                  </div>
                </ScrollArea>

                {/* Input Area */}
                <div className="flex gap-2">
                  <Textarea
                    ref={inputRef}
                    placeholder="Ask about project management, platform features, or get general guidance..."
                    value={inputMessage}
                    onChange={(e) => setInputMessage(e.target.value)}
                    onKeyPress={handleKeyPress}
                    disabled={isLoading}
                    className="flex-1 min-h-[40px] max-h-32 resize-none"
                    rows={1}
                  />
                  <Button
                    onClick={handleSendMessage}
                    disabled={!inputMessage.trim() || isLoading}
                    size="sm"
                    className="self-end"
                  >
                    {isLoading ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Send className="h-4 w-4" />
                    )}
                  </Button>
                </div>
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center bg-card">
                <div className="text-center">
                  <Bot className="h-16 w-16 mx-auto mb-4 text-muted-foreground/50" />
                  <h3 className="text-lg font-semibold mb-2">Welcome to ProjectNest AI</h3>
                  <p className="text-muted-foreground mb-4 max-w-md">
                    Create a new conversation or select an existing one to start getting help with your projects and platform usage.
                  </p>
                  <Button onClick={() => setShowCreateChatDialog(true)}>
                    <Plus className="h-4 w-4 mr-2" />
                    Start New Conversation
                  </Button>
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Edit Conversation Dialog */}
      <Dialog open={!!editingConversation} onOpenChange={(open) => !open && setEditingConversation(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Rename Conversation</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Input
              placeholder="Enter new name..."
              value={editingChatName}
              onChange={(e) => setEditingChatName(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleUpdateConversation()}
            />
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => setEditingConversation(null)}>
                Cancel
              </Button>
              <Button onClick={handleUpdateConversation} disabled={!editingChatName.trim()}>
                Update
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};