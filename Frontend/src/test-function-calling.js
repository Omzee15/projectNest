// Test file for Gemini Function Calling Implementation
// This file demonstrates proper function calling with Google Gemini AI

import { GoogleGenerativeAI } from '@google/generative-ai';

// Initialize Gemini AI
const apiKey = process.env.VITE_GEMINI_API_KEY;
if (!apiKey) {
  throw new Error('VITE_GEMINI_API_KEY is required');
}
const genAI = new GoogleGenerativeAI(apiKey);

// Define the project creation function declaration
const createProjectFunctionDeclaration = {
  name: 'create_project',
  description: 'Creates a new project with tasks, lists, and additional components like flowcharts and database schemas. Use this when the user asks to create, build, or set up a project.',
  parameters: {
    type: 'object',
    properties: {
      project_name: {
        type: 'string',
        description: 'A clear, concise name for the project'
      },
      project_description: {
        type: 'string',
        description: 'A brief description of what the project is about'
      },
      lists: {
        type: 'array',
        description: 'Array of task lists/categories for organizing project tasks',
        items: {
          type: 'object',
          properties: {
            name: {
              type: 'string',
              description: 'Name of the task list (e.g., "Planning", "Development", "Testing")'
            },
            tasks: {
              type: 'array',
              description: 'Tasks within this list',
              items: {
                type: 'object',
                properties: {
                  title: {
                    type: 'string',
                    description: 'Task title'
                  },
                  description: {
                    type: 'string',
                    description: 'Detailed task description'
                  }
                },
                required: ['title']
              }
            }
          },
          required: ['name', 'tasks']
        }
      },
      flowchart: {
        type: 'string',
        description: 'Mermaid flowchart script if the project needs a visual workflow'
      },
      database_schema: {
        type: 'string',
        description: 'Database schema description if the project involves data storage'
      }
    },
    required: ['project_name', 'project_description', 'lists']
  }
};

// Mock function to simulate project creation
function createProject(projectData) {
  console.log('Creating project with data:', JSON.stringify(projectData, null, 2));
  return {
    success: true,
    project_id: 'proj_' + Date.now(),
    message: `Project "${projectData.project_name}" created successfully`
  };
}

// Test function calling implementation
export async function testFunctionCalling() {
  try {
    const model = genAI.getGenerativeModel({ 
      model: 'gemini-2.5-flash',
      tools: [{
        functionDeclarations: [createProjectFunctionDeclaration]
      }]
    });

    // Test prompt that should trigger function calling
    const testPrompt = `I want to create a project for building a task management app called "TaskMaster". 

    The app should help users:
    - Create and organize tasks
    - Set priorities and due dates
    - Track progress
    - Collaborate with team members
    
    I want to organize this into different phases like Planning, Development, Testing, and Deployment. Each phase should have specific tasks.
    
    Please create a complete project structure for this.`;

    console.log('🚀 Testing Function Calling...\n');
    console.log('📝 Test Prompt:', testPrompt);
    console.log('\n⏳ Sending request to Gemini...\n');

    // Send request with function declarations
    const result = await model.generateContent(testPrompt);
    const response = result.response;

    console.log('📨 Response received!');
    console.log('🔍 Checking for function calls...\n');

    // Check if there are function calls
    const functionCalls = response.functionCalls?.();
    
    if (functionCalls && functionCalls.length > 0) {
      console.log('✅ Function call detected!');
      console.log('📞 Function calls:', functionCalls.length);
      
      for (const functionCall of functionCalls) {
        console.log(`\n🔧 Function: ${functionCall.name}`);
        console.log('📋 Arguments:', JSON.stringify(functionCall.args, null, 2));
        
        if (functionCall.name === 'create_project') {
          // Execute the actual function
          const result = createProject(functionCall.args);
          console.log('✅ Function execution result:', result);
          
          // Send the result back to the model for a user-friendly response
          const finalResult = await model.generateContent([
            {
              role: 'user',
              parts: [{ text: testPrompt }]
            },
            {
              role: 'model', 
              parts: response.candidates[0].content.parts
            },
            {
              role: 'user',
              parts: [{
                functionResponse: {
                  name: functionCall.name,
                  response: result
                }
              }]
            }
          ]);
          
          console.log('\n💬 Final AI Response:');
          console.log(finalResult.response.text());
        }
      }
    } else {
      console.log('❌ No function calls found');
      console.log('📄 Direct response:', response.text());
    }

  } catch (error) {
    console.error('❌ Error in function calling test:', error);
  }
}

// Additional test with explicit confirmation flow
export async function testConfirmationFlow() {
  try {
    const model = genAI.getGenerativeModel({ 
      model: 'gemini-2.5-flash',
      tools: [{
        functionDeclarations: [createProjectFunctionDeclaration]
      }]
    });

    console.log('\n🔄 Testing Confirmation Flow...\n');

    // First message - project idea
    const conversation = [
      {
        role: 'user',
        parts: [{ text: 'I want to build an e-commerce website for selling handmade crafts' }]
      }
    ];

    let result = await model.generateContent(conversation);
    let response = result.response.text();
    console.log('🤖 AI Response 1:', response);

    // Add AI response to conversation
    conversation.push({
      role: 'model',
      parts: [{ text: response }]
    });

    // User confirms project creation
    conversation.push({
      role: 'user',
      parts: [{ text: 'Yes, please go ahead and create the project structure for me.' }]
    });

    console.log('\n👤 User confirms: "Yes, please go ahead and create the project structure for me."');
    console.log('\n⏳ Sending confirmation...\n');

    // Send confirmation
    result = await model.generateContent(conversation);
    response = result.response;

    const functionCalls = response.functionCalls?.();
    
    if (functionCalls && functionCalls.length > 0) {
      console.log('✅ Function call triggered by confirmation!');
      
      for (const functionCall of functionCalls) {
        console.log(`\n🔧 Function: ${functionCall.name}`);
        console.log('📋 Arguments:', JSON.stringify(functionCall.args, null, 2));
        
        if (functionCall.name === 'create_project') {
          const result = createProject(functionCall.args);
          console.log('✅ Project created:', result);
        }
      }
    } else {
      console.log('❌ No function call after confirmation');
      console.log('📄 Response:', response.text());
    }

  } catch (error) {
    console.error('❌ Error in confirmation flow test:', error);
  }
}

// Run tests if this file is executed directly
if (typeof window === 'undefined') {
  // Node.js environment
  testFunctionCalling().then(() => {
    return testConfirmationFlow();
  });
}