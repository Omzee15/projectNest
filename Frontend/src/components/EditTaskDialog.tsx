import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ColorPicker } from '@/components/ui/color-picker';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Loader2, Users } from 'lucide-react';
import { Task, TaskUpdateRequest } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { TaskAssigneeSelector } from './TaskAssigneeSelector';

interface EditTaskDialogProps {
  task?: Task;
  projectId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onTaskUpdated: (updatedTask: Task) => void;
}

export function EditTaskDialog({ task, projectId, open, onOpenChange, onTaskUpdated }: EditTaskDialogProps) {
  const [formData, setFormData] = useState<TaskUpdateRequest>({
    title: '',
    description: '',
    priority: undefined,
    status: '',
    color: '#FFFFFF',
    due_date: '',
  });
  const [assignedUserIds, setAssignedUserIds] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  // Reset form when task changes
  useEffect(() => {
    if (task && open) {
      setFormData({
        title: task.title,
        description: task.description || '',
        priority: task.priority || undefined,
        status: task.status,
        color: task.color || '#FFFFFF',
        due_date: task.due_date ? task.due_date.split('T')[0] : '',
      });
      setAssignedUserIds(task.assignees?.map((a) => a.user_uid) || []);
    }
  }, [task, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!task) return;

    try {
      setIsLoading(true);

      // Filter out empty values and prepare the update request
      const updates: TaskUpdateRequest = {};
      if (formData.title?.trim() && formData.title !== task.title) {
        updates.title = formData.title.trim();
      }
      if (formData.description !== task.description) {
        updates.description = formData.description || undefined;
      }
      if (formData.priority !== task.priority) {
        updates.priority = formData.priority;
      }
      if (formData.status && formData.status !== task.status) {
        updates.status = formData.status;
      }
      if (formData.color && formData.color !== task.color) {
        updates.color = formData.color;
      }
      if (formData.due_date !== (task.due_date ? task.due_date.split('T')[0] : '')) {
        updates.due_date = formData.due_date || undefined;
      }

      // Check for assignee changes
      const originalAssigneeIds = task.assignees?.map((a) => a.user_uid) || [];
      const assigneeIdsChanged = 
        assignedUserIds.length !== originalAssigneeIds.length ||
        assignedUserIds.some((id) => !originalAssigneeIds.includes(id)) ||
        originalAssigneeIds.some((id) => !assignedUserIds.includes(id));

      // Only send request if there are changes
      if (Object.keys(updates).length === 0 && !assigneeIdsChanged) {
        toast({
          title: 'No changes',
          description: 'No changes were made to the task.',
        });
        onOpenChange(false);
        return;
      }

      // Update task fields if there are any
      let updatedTask = task;
      if (Object.keys(updates).length > 0) {
        const response = await apiService.partialUpdateTask(task.task_uid, updates);
        updatedTask = response.data;
      }

      // Update assignees if changed
      if (assigneeIdsChanged) {
        await apiService.bulkAssignTask(task.task_uid, assignedUserIds);
        // Fetch updated task with assignees
        const taskResponse = await apiService.getProject(task.list_id.toString());
        // Find the updated task in the response
        // For now, we'll just add assignees to the response
        updatedTask = { ...updatedTask, assignees: undefined }; // Backend should return this
      }
      
      onTaskUpdated(updatedTask);
      onOpenChange(false);
      
      toast({
        title: 'Task updated',
        description: 'Task has been successfully updated.',
      });
    } catch (error) {
      console.error('Failed to update task:', error);
      toast({
        title: 'Error',
        description: 'Failed to update task. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (field: keyof TaskUpdateRequest, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Task</DialogTitle>
          <DialogDescription>
            Update the task details below.
          </DialogDescription>
        </DialogHeader>

        {/* Display Current Assignees */}
        {task?.assignees && task.assignees.length > 0 && (
          <div className="bg-muted/50 rounded-lg p-4 space-y-3 border border-border">
            <div className="flex items-center gap-2 text-sm font-medium text-foreground">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span>Assigned Members ({task.assignees.length})</span>
            </div>
            <div className="flex flex-wrap gap-2">
              {task.assignees.map((assignee) => (
                <Badge
                  key={assignee.user_uid}
                  variant="secondary"
                  className="flex items-center gap-2 py-1.5 px-3"
                >
                  <Avatar className="h-5 w-5 border border-background">
                    <AvatarFallback className="text-[10px] bg-primary/10 text-primary">
                      {getInitials(assignee.name)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-xs font-medium">{assignee.name}</span>
                </Badge>
              ))}
            </div>
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              value={formData.title}
              onChange={(e) => handleInputChange('title', e.target.value)}
              placeholder="Enter task title"
              required
              className="w-full"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={formData.description}
              onChange={(e) => handleInputChange('description', e.target.value)}
              placeholder="Enter task description (optional)"
              rows={3}
              className="w-full resize-none"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="priority">Priority</Label>
              <Select
                value={formData.priority || ''}
                onValueChange={(value) => handleInputChange('priority', value || undefined)}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select priority" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="status">Status</Label>
              <Select
                value={formData.status}
                onValueChange={(value) => handleInputChange('status', value)}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="todo">To Do</SelectItem>
                  <SelectItem value="in_progress">In Progress</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Color</Label>
            <div className="w-full">
              <ColorPicker
                value={formData.color || '#FFFFFF'}
                onChange={(color) => handleInputChange('color', color)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label>Assign To</Label>
            <TaskAssigneeSelector
              projectId={projectId}
              selectedUserIds={assignedUserIds}
              onSelectionChange={setAssignedUserIds}
              disabled={isLoading}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="due_date">Due Date</Label>
            <Input
              id="due_date"
              type="date"
              value={formData.due_date}
              onChange={(e) => handleInputChange('due_date', e.target.value)}
              className="w-full"
            />
          </div>

          <DialogFooter className="flex flex-col sm:flex-row gap-2 sm:gap-0">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
              className="w-full sm:w-auto"
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading} className="w-full sm:w-auto">
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Updating...
                </>
              ) : (
                'Update Task'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}