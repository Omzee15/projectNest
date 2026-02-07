import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Calendar, Flag, MoreVertical, Trash2, Loader2, Edit, CheckCircle } from 'lucide-react';
import { Task } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { ColorIndicator } from '@/components/ui/color-picker';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { getCurrentUserUid } from '@/utils/auth';
import { useState } from 'react';

interface TaskCardProps {
  task: Task;
  onClick?: () => void;
  onTaskDelete?: (taskUid: string) => void;
  onTaskUpdate?: (taskUid: string, updates: Partial<Task>) => void;
  onTaskEdit?: (task: Task) => void;
}

export function TaskCard({ task, onClick, onTaskDelete, onTaskUpdate, onTaskEdit }: TaskCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: task.task_uid });
  
  const [isDeleting, setIsDeleting] = useState(false);
  const [isUpdatingCompletion, setIsUpdatingCompletion] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const { toast } = useToast();

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const handleDeleteTask = async () => {
    try {
      setIsDeleting(true);
      await onTaskDelete?.(task.task_uid);
      setShowDeleteDialog(false);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete task. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const handleToggleCompletion = async (checked: boolean) => {
    try {
      setIsUpdatingCompletion(true);
      
      // Call API to update task completion using partial update
      const response = await apiService.partialUpdateTask(task.task_uid, { is_completed: checked });
      
      // Update local state through parent callback
      if (onTaskUpdate) {
        onTaskUpdate(task.task_uid, response.data);
      }
      
      toast({
        title: checked ? 'Task completed!' : 'Task uncompleted',
        description: checked ? `"${task.title}" has been marked as completed.` : `"${task.title}" has been marked as incomplete.`,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to update task completion. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsUpdatingCompletion(false);
    }
  };

  const handleEditTask = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (onTaskEdit) {
      onTaskEdit(task);
    }
  };

  const getPriorityColor = (priority?: string) => {
    switch (priority) {
      case 'high':
        return 'destructive';
      case 'medium':
        return 'warning';
      case 'low':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  const getStatusBorderClass = (status?: string) => {
    switch (status) {
      case 'completed':
        return 'border-l-4 border-l-green-500';
      case 'in_progress':
        return 'border-l-4 border-l-yellow-500';
      case 'todo':
        return 'border-l-4 border-l-gray-300';
      default:
        return 'border-l-4 border-l-gray-300';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    });
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const currentUserUid = getCurrentUserUid();
  const isAssignedToMe = task.assignees?.some((assignee) => {
    const match = assignee.user_uid === currentUserUid;
    return match;
  });
  
  // Debug logging
  console.log('=== TaskCard Debug ===');
  console.log('Task:', task.title);
  console.log('Current User UID:', currentUserUid, 'Type:', typeof currentUserUid);
  console.log('Task Assignees:', task.assignees);
  if (task.assignees && task.assignees.length > 0) {
    task.assignees.forEach(assignee => {
      console.log('  Assignee UID:', assignee.user_uid, 'Type:', typeof assignee.user_uid, 'Match:', assignee.user_uid === currentUserUid);
    });
  }
  console.log('Is Assigned to Me:', isAssignedToMe);
  console.log('=====================');

  return (
    <Card
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`
        cursor-pointer transition-all duration-200 hover:shadow-md group relative
        ${isDragging ? 'opacity-50 rotate-2 scale-105' : ''}
        ${task.is_completed ? 'opacity-75 bg-muted' : 'bg-card hover:bg-card-hover'}
        ${getStatusBorderClass(task.status)}
        ${isAssignedToMe ? 'ring-1 ring-primary/50' : ''}
        border-border-light
      `}
      onClick={onClick}
    >
      <CardContent className="p-3 sm:p-4 space-y-2 sm:space-y-3">
        <div className="space-y-2">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-2 flex-1 min-w-0">
              <div 
                onClick={(e) => e.stopPropagation()}
                className="flex items-center gap-2 flex-shrink-0"
              >
                <Checkbox
                  checked={task.is_completed}
                  onCheckedChange={handleToggleCompletion}
                  disabled={isUpdatingCompletion}
                  className="flex-shrink-0"
                />
                <ColorIndicator color={task.color || '#FFFFFF'} size="sm" />
              </div>
              <h4 className={`text-sm font-medium leading-tight flex-1 break-words word-wrap min-w-0 ${
                task.is_completed ? 'line-through text-muted-foreground' : 'text-card-foreground'
              }`}>
                {task.title}
              </h4>
              {task.is_completed && (
                <CheckCircle className="h-4 w-4 text-green-500 flex-shrink-0" />
              )}
            </div>
            
            {/* Delete Button */}
            <div className="opacity-0 group-hover:opacity-100 transition-opacity">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button 
                    variant="ghost" 
                    size="sm"
                    className="h-6 w-6 p-0"
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                  >
                    <MoreVertical className="h-3 w-3" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem
                    className="cursor-pointer"
                    onClick={handleEditTask}
                  >
                    <Edit className="mr-2 h-4 w-4" />
                    Edit Task
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="text-destructive focus:text-destructive cursor-pointer"
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      setShowDeleteDialog(true);
                    }}
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete Task
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              
              <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete Task</AlertDialogTitle>
                    <AlertDialogDescription>
                      Are you sure you want to delete "{task.title}"? This action cannot be undone.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleDeleteTask}
                      className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                      disabled={isDeleting}
                    >
                      {isDeleting ? (
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
            </div>
          </div>
          
          {task.description && (
            <p className="text-xs text-muted-foreground line-clamp-2">
              {task.description}
            </p>
          )}
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap">
            {task.priority && (
              <Badge variant="secondary" className="text-xs px-2 py-0.5">
                <Flag className="w-3 h-3 mr-1" />
                {task.priority}
              </Badge>
            )}
            {task.categories && task.categories.length > 0 && (
              <div className="flex items-center gap-1">
                {task.categories.length === 1 ? (
                  <Badge
                    variant="outline"
                    className="text-xs px-2 py-0.5"
                    style={{ 
                      backgroundColor: task.categories[0].color + '20',
                      borderColor: task.categories[0].color,
                      color: task.categories[0].color
                    }}
                  >
                    {task.categories[0].name}
                  </Badge>
                ) : (
                  task.categories.slice(0, 3).map((category) => (
                    <div
                      key={category.category_uid}
                      className="w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-semibold border-2 border-background"
                      style={{ backgroundColor: category.color, color: 'white' }}
                      title={category.name}
                    >
                      {category.name.substring(0, 2).toUpperCase()}
                    </div>
                  ))
                )}
                {task.categories.length > 3 && (
                  <div className="w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-semibold bg-muted text-muted-foreground border-2 border-background">
                    +{task.categories.length - 3}
                  </div>
                )}
              </div>
            )}
            {task.assignees && task.assignees.length > 0 && (
              <div className="flex -space-x-2">
                {task.assignees.slice(0, 3).map((assignee) => (
                  <Avatar 
                    key={assignee.user_uid} 
                    className="h-6 w-6 border-2 border-background"
                    title={assignee.name}
                  >
                    <AvatarFallback className="text-[10px] bg-primary/10 text-primary">
                      {getInitials(assignee.name)}
                    </AvatarFallback>
                  </Avatar>
                ))}
                {task.assignees.length > 3 && (
                  <Avatar className="h-6 w-6 border-2 border-background">
                    <AvatarFallback className="text-[10px] bg-muted text-muted-foreground">
                      +{task.assignees.length - 3}
                    </AvatarFallback>
                  </Avatar>
                )}
              </div>
            )}
          </div>

          {task.due_date && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Calendar className="w-3 h-3" />
              <span>{formatDate(task.due_date)}</span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}