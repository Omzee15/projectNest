import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { 
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Navbar } from '@/components/Navbar';
import { CreateProjectDialog } from '@/components/CreateProjectDialog';
import { AnimatedProgress } from '@/components/AnimatedProgress';
import { DashboardAIChat } from '@/components/DashboardAIChat';
import { Project, ProjectWithProgress } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { Plus, Calendar, Users, ArrowRight, Loader2, Trash2, MoreVertical, Bot } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

export default function Projects() {
  const [projects, setProjects] = useState<ProjectWithProgress[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [deletingProject, setDeletingProject] = useState<string | null>(null);
  const [projectToDelete, setProjectToDelete] = useState<{uid: string, name: string} | null>(null);
  const [isAIChatOpen, setIsAIChatOpen] = useState(false);
  const [isDevAISelectOpen, setIsDevAISelectOpen] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      setIsLoading(true);
      const response = await apiService.getProjectsWithProgress();
      setProjects(response.data || []);
    } catch (error) {
      console.error('Failed to load projects:', error);
      toast({
        title: 'Error',
        description: 'Failed to load projects. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

    const handleProjectCreated = (project: Project) => {
    loadProjects();
    toast({
      title: 'Success',
      description: 'Project created successfully.',
    });
  };

  const navigateToProjectDevAI = (projectUid: string) => {
    setIsDevAISelectOpen(false);
    navigate(`/project/${projectUid}/dev-ai`);
  };

    const handleDeleteProject = async (projectUid: string, projectName: string) => {
    try {
      setDeletingProject(projectUid);
      await apiService.deleteProject(projectUid);
      setProjects(prev => prev.filter(p => p.project_uid !== projectUid));
      setProjectToDelete(null);
      toast({
        title: 'Project deleted',
        description: `Project "${projectName}" has been deleted`,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete project. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setDeletingProject(null);
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'No date set';
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'completed':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'on_hold':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="container mx-auto px-6 py-16">
          <div className="flex items-center justify-center min-h-[400px]">
            <div className="flex flex-col items-center gap-4">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-muted-foreground">Loading projects...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      
      <main className="w-full max-w-[90vw] mx-auto px-4 sm:px-6 py-6 sm:py-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6 sm:mb-8">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-foreground">Projects</h1>
            <p className="text-sm sm:text-base text-muted-foreground mt-2">
              Manage and organize all your projects in one place
            </p>
          </div>
          <CreateProjectDialog 
            onProjectCreate={handleProjectCreated}
            trigger={
              <Button size="lg" className="w-full sm:w-auto">
                <Plus className="mr-2 h-4 w-4" />
                New Project
              </Button>
            }
          />
        </div>

        {/* Projects Grid */}
        {projects.length === 0 ? (
          <div className="text-center py-12 sm:py-16 px-4">
            <div className="max-w-md mx-auto">
              <div className="mb-6">
                <div className="w-16 h-16 bg-muted rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Users className="h-8 w-8 text-muted-foreground" />
                </div>
                <h3 className="text-lg sm:text-xl font-semibold text-foreground mb-2">
                  No projects yet
                </h3>
                <p className="text-sm sm:text-base text-muted-foreground mb-6">
                  Get started by creating your first project. Organize tasks, collaborate with your team, and track progress.
                </p>
                <CreateProjectDialog 
                  onProjectCreate={handleProjectCreated}
                  trigger={
                    <Button size="lg" className="w-full sm:w-auto">
                      <Plus className="mr-2 h-4 w-4" />
                      Create Your First Project
                    </Button>
                  }
                />
              </div>
            </div>
          </div>
        ) : (
          <div className="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
            {projects.map((project) => (
                <Card 
                  key={project.project_uid} 
                  className="hover:shadow-lg transition-all duration-200 group relative"
                >
                <div
                  className="cursor-pointer"
                  onClick={() => navigate(`/project/${project.project_uid}`)}
                >
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        <CardTitle className="text-lg sm:text-xl group-hover:text-primary transition-colors truncate">
                          {project.name}
                        </CardTitle>
                        <div className="flex items-center gap-2 mt-2 flex-wrap">
                          <span 
                            className={`px-2 py-1 text-xs rounded-md border ${getStatusColor(project.status)}`}
                          >
                            {project.status.replace('_', ' ').toUpperCase()}
                          </span>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <ArrowRight className="h-5 w-5 text-muted-foreground group-hover:text-primary group-hover:translate-x-1 transition-all" />
                      </div>
                    </div>
                    {project.description && (
                      <CardDescription className="mt-2 line-clamp-3">
                        {project.description}
                      </CardDescription>
                    )}
                  </CardHeader>
                  <CardContent>
                    {/* Progress Section */}
                    <div className="mb-4">
                      <div className="flex items-center justify-between text-sm text-muted-foreground mb-2">
                        <span>Progress</span>
                        <span>{project.task_stats?.completed_tasks || 0}/{project.task_stats?.total_tasks || 0} tasks</span>
                      </div>
                      <AnimatedProgress 
                        value={(project.task_stats?.progress || 0) * 100} 
                        className="h-2" 
                      />
                      
                      <div className="flex gap-4 mt-3 text-xs">
                        <span className="flex items-center gap-1">
                          <div className="w-2 h-2 bg-gray-500 rounded-full"></div>
                          {project.task_stats?.todo_tasks || 0} To Do
                        </span>
                        <span className="flex items-center gap-1">
                          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                          {project.task_stats?.completed_tasks || 0} Done
                        </span>
                      </div>
                    </div>
                    
                    <div className="space-y-2 text-sm text-muted-foreground">
                      <div className="flex items-center gap-2">
                        <Calendar className="h-4 w-4" />
                        <span>Created: {formatDate(project.created_at)}</span>
                      </div>
                      {project.start_date && (
                        <div className="flex items-center gap-2">
                          <Calendar className="h-4 w-4" />
                          <span>Start: {formatDate(project.start_date)}</span>
                        </div>
                      )}
                      {project.end_date && (
                        <div className="flex items-center gap-2">
                          <Calendar className="h-4 w-4" />
                          <span>End: {formatDate(project.end_date)}</span>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </div>
                
                {/* Actions Dropdown */}
                <div className="absolute top-4 right-4">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button 
                        variant="ghost" 
                        size="sm"
                        className="opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8 p-0"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem
                        className="text-destructive focus:text-destructive cursor-pointer"
                        onClick={(e) => {
                          e.stopPropagation();
                          setProjectToDelete({uid: project.project_uid, name: project.name});
                        }}
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Delete Project
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </Card>
            ))}
          </div>
        )}
      </main>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!projectToDelete} onOpenChange={(open) => !open && setProjectToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Project</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{projectToDelete?.name}"? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => projectToDelete && handleDeleteProject(projectToDelete.uid, projectToDelete.name)}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={!!deletingProject}
            >
              {deletingProject ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete'
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* AI Chat Floating Buttons */}
      <div className="fixed bottom-4 sm:bottom-6 right-4 sm:right-6 z-40 flex gap-2 sm:gap-4">
        <Button
          size="lg"
          onClick={() => setIsDevAISelectOpen(true)}
          className="rounded-full h-12 w-12 sm:h-14 sm:w-14 shadow-lg hover:shadow-xl transition-shadow bg-muted p-0"
          title="Open Project Dev AI"
        >
          <Bot className="h-5 w-5 sm:h-6 sm:w-6 text-blue-500" />
        </Button>
        <Button
          size="lg"
          onClick={() => setIsAIChatOpen(true)}
          className="rounded-full h-12 w-12 sm:h-14 sm:w-14 shadow-lg hover:shadow-xl transition-shadow p-0"
          title="Open AI Assistant"
        >
          <Bot className="h-5 w-5 sm:h-6 sm:w-6" />
        </Button>
      </div>
      {/* AI Chat Component */}
      <DashboardAIChat 
        isOpen={isAIChatOpen} 
        onClose={() => setIsAIChatOpen(false)} 
      />

      {/* Dev AI Project Selection Dialog */}
      <Dialog open={isDevAISelectOpen} onOpenChange={setIsDevAISelectOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Select Project to Open Dev AI</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {projects.map((project) => (
              <Button
                key={project.project_uid}
                variant="outline"
                className="flex items-center justify-between w-full p-4"
                onClick={() => navigateToProjectDevAI(project.project_uid)}
              >
                <span className="font-medium">{project.name}</span>
                <span className="text-sm text-muted-foreground">
                  Open Dev AI
                </span>
              </Button>
            ))}
            {projects.length === 0 && (
              <p className="text-center text-muted-foreground py-4">
                No projects found. Create a project first to use Dev AI.
              </p>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}