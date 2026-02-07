import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ProjectBoard } from '@/components/ProjectBoard';
import { Navbar } from '@/components/Navbar';
import { TaskDetailsDialog } from '@/components/TaskDetailsDialog';
import { EditTaskDialog } from '@/components/EditTaskDialog';
import { EnhancedDevAIDialog } from '@/components/EnhancedDevAIDialog';
import { ManageCategoriesDialog } from '@/components/ManageCategoriesDialog';
import { ProjectWithLists, Project, ListWithTasks, Task } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';



export default function ProjectDashboard() {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const [project, setProject] = useState<ProjectWithLists | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [showTaskDetailsDialog, setShowTaskDetailsDialog] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [showEditTaskDialog, setShowEditTaskDialog] = useState(false);
  const [showManageCategoriesDialog, setShowManageCategoriesDialog] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    loadProject();
  }, [projectId]);

  const loadProject = async () => {
    if (!projectId) {
      toast({
        title: 'No Project Selected',
        description: 'Redirecting to projects list...',
        variant: 'default',
      });
      navigate('/projects');
      return;
    }

    try {
      setIsLoading(true);
      const response = await apiService.getProject(projectId);
      setProject(response.data);
    } catch (error) {
      console.error('Failed to load project:', error);
      
      // Check if this is an authentication error
      if (error instanceof Error && error.message.includes('Authentication failed')) {
        // Don't handle auth errors here, let them propagate to AuthGuard
        throw error;
      }
      
      toast({
        title: 'Project Not Found',
        description: 'This project does not exist or you do not have access to it.',
        variant: 'destructive',
      });
      
      // Redirect to projects list
      navigate('/projects');
    } finally {
      setIsLoading(false);
    }
  };

  const handleTaskMove = async (taskId: string, newListUid: string) => {
    try {
      await apiService.moveTask(taskId, newListUid);
      toast({
        title: 'Task moved',
        description: 'Task successfully moved to new list.',
      });
      // Reload the project to reflect changes
      await loadProject();
    } catch (error) {
      console.error('Failed to move task:', error);
      toast({
        title: 'Move failed',
        description: 'Could not move task. Please try again.',
        variant: 'destructive',
      });
    }
  };

  const handleTaskClick = (taskId: string) => {
    console.log('Task clicked:', taskId);
    
    // Find the task in the project lists
    if (project) {
      for (const list of project.lists) {
        const task = (list.tasks || []).find(t => t.task_uid === taskId);
        if (task) {
          setSelectedTask(task);
          setShowTaskDetailsDialog(true);
          break;
        }
      }
    }
  };

  const handleEditTaskFromDetails = (task: Task) => {
    setShowTaskDetailsDialog(false);
    setEditingTask(task);
    setShowEditTaskDialog(true);
  };

  const handleTaskUpdated = (updatedTask: Task) => {
    // Update the task in the project state
    handleTaskUpdate(updatedTask.task_uid, updatedTask);
    
    // Update the selected task if it's the same one
    if (selectedTask && selectedTask.task_uid === updatedTask.task_uid) {
      setSelectedTask(updatedTask);
    }
    
    // Close the edit dialog
    setShowEditTaskDialog(false);
    setEditingTask(null);
  };

  const handleAddTask = async (listUid: string, taskData: any) => {
    try {
      // taskData already has the correct structure from CreateTaskDialog
      await apiService.createTask(taskData);
      
      toast({
        title: 'Task created',
        description: `Task "${taskData.title}" created successfully`,
      });

      // Reload the project to reflect changes
      await loadProject();
    } catch (error) {
      console.error('Failed to create task:', error);
      toast({
        title: 'Failed to create task',
        description: error instanceof Error ? error.message : 'An unexpected error occurred',
        variant: 'destructive',
      });
    }
  };

  const handleAddList = async () => {
    console.log('List created, reloading project data');
    // Reload the project to get the updated lists
    await loadProject();
  };

  const handleDeleteList = async (listUid: string) => {
    try {
      await apiService.deleteList(listUid);
      
      toast({
        title: 'List deleted',
        description: 'List deleted successfully',
      });

      // Update local state to remove the deleted list
      if (project) {
        setProject({
          ...project,
          lists: project.lists.filter(list => list.list_uid !== listUid)
        });
      }
    } catch (error) {
      console.error('Failed to delete list:', error);
      toast({
        title: 'Failed to delete list',
        description: error instanceof Error ? error.message : 'An unexpected error occurred',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteTask = async (taskUid: string) => {
    try {
      await apiService.deleteTask(taskUid);
      
      toast({
        title: 'Task deleted',
        description: 'Task deleted successfully',
      });

      // Update local state to remove the deleted task
      if (project) {
        setProject({
          ...project,
          lists: project.lists.map(list => ({
            ...list,
            tasks: list.tasks.filter(task => task.task_uid !== taskUid)
          }))
        });
      }
    } catch (error) {
      console.error('Failed to delete task:', error);
      toast({
        title: 'Failed to delete task',
        description: error instanceof Error ? error.message : 'An unexpected error occurred',
        variant: 'destructive',
      });
    }
  };

  const handleProjectUpdate = (updatedProject: Project) => {
    if (project) {
      setProject({
        ...project,
        ...updatedProject,
      });
    }
  };

  const handleListUpdate = (listUid: string, updatedList: ListWithTasks) => {
    if (project) {
      setProject({
        ...project,
        lists: project.lists.map(list =>
          list.list_uid === listUid ? updatedList : list
        )
      });
    }
  };

  const handleTaskUpdate = (taskUid: string, updatedTask: Task) => {
    if (project) {
      const updatedLists = project.lists.map(list => ({
        ...list,
        tasks: (list.tasks || []).map(task => 
          task.task_uid === taskUid ? updatedTask : task
        )
      }));

      setProject({
        ...project,
        lists: updatedLists
      });
    }
  };

  const handleListMove = async (listId: string, newPosition: number) => {
    // Update local state immediately for responsive UI
    if (project) {
      const updatedLists = [...project.lists];
      const listIndex = updatedLists.findIndex(list => list.list_uid === listId);
      
      if (listIndex !== -1) {
        const [movedList] = updatedLists.splice(listIndex, 1);
        updatedLists.splice(newPosition, 0, movedList);
        
        setProject({
          ...project,
          lists: updatedLists
        });
      }
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading project...</p>
          </div>
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-foreground mb-2">Project not found</h2>
            <p className="text-muted-foreground">The project you're looking for doesn't exist.</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen bg-background flex flex-col overflow-hidden">
      <Navbar />
      <ProjectBoard
        project={project}
        onTaskMove={handleTaskMove}
        onTaskClick={handleTaskClick}
        onAddTask={handleAddTask}
        onAddList={handleAddList}
        onDeleteList={handleDeleteList}
        onDeleteTask={handleDeleteTask}
        onProjectUpdate={handleProjectUpdate}
        onListUpdate={handleListUpdate}
        onTaskUpdate={handleTaskUpdate}
        onListMove={handleListMove}
        onManageCategories={() => setShowManageCategoriesDialog(true)}
      />
      
      <TaskDetailsDialog
        task={selectedTask || undefined}
        open={showTaskDetailsDialog}
        onOpenChange={setShowTaskDetailsDialog}
        onEditTask={handleEditTaskFromDetails}
        projectUid={project?.project_uid}
      />
      
      <EditTaskDialog
        task={editingTask || undefined}
        projectId={projectId || ''}
        open={showEditTaskDialog}
        onOpenChange={(open) => {
          setShowEditTaskDialog(open);
          if (!open) setEditingTask(null);
        }}
        onTaskUpdated={handleTaskUpdated}
      />

      <ManageCategoriesDialog
        projectUid={project.project_uid}
        open={showManageCategoriesDialog}
        onOpenChange={setShowManageCategoriesDialog}
        onCategoriesUpdated={loadProject}
      />

      {/* Enhanced AI Floating Button - Can create/edit lists and tasks */}
      {project && <EnhancedDevAIDialog projectUid={project.project_uid} />}
    </div>
  );
}