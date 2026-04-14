import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { ColorPicker } from '@/components/ui/color-picker';
import { Plus, Calendar, Flag, Lock } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { COLORS } from '@/types';
import { TaskAssigneeSelector } from './TaskAssigneeSelector';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

interface CreateTaskDialogProps {
  listUid: string;
  listName: string;
  projectId: string;
  onTaskCreate?: (taskData: any) => void;
  trigger?: React.ReactNode;
  canWrite?: boolean;
}

export function CreateTaskDialog({ listUid, listName, projectId, onTaskCreate, trigger, canWrite = true }: CreateTaskDialogProps) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState('medium');
  const [color, setColor] = useState<string>(COLORS.WHITE);
  const [dueDate, setDueDate] = useState('');
  const [assignedUserIds, setAssignedUserIds] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) {
      toast({
        title: 'Error',
        description: 'Task title is required',
        variant: 'destructive',
      });
      return;
    }

    setIsLoading(true);

    try {
      const taskData = {
        list_uid: listUid,
        title: title.trim(),
        description: description.trim() || undefined,
        priority: priority || undefined,
        color,
        status: 'todo',
        due_date: dueDate || undefined,
        assigned_user_ids: assignedUserIds.length > 0 ? assignedUserIds : undefined,
      };

      // Call callback with the task data - parent will handle API call
      await onTaskCreate?.(taskData);
      
      // Reset form
      setTitle('');
      setDescription('');
      setPriority('medium');
      setColor(COLORS.WHITE);
      setDueDate('');
      setAssignedUserIds([]);
      setOpen(false);

      toast({
        title: 'Success',
        description: `Task "${title}" has been created successfully`,
      });
    } catch (error) {
      console.error('Failed to create task:', error);
      toast({
        title: 'Error',
        description: 'Failed to create task. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(newOpen) => {
      if (!canWrite && newOpen) {
        toast({
          title: 'Access Denied',
          description: 'You do not have permission to create tasks. Only users with edit access can create tasks.',
          variant: 'destructive',
        });
        return;
      }
      setOpen(newOpen);
    }}>
      <DialogTrigger asChild>
        {trigger || (
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full justify-start text-muted-foreground hover:text-foreground hover:bg-secondary"
                  disabled={!canWrite}
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add a card
                </Button>
              </TooltipTrigger>
              {!canWrite && (
                <TooltipContent>
                  <p>You have view-only access. Only users with edit permission can create tasks.</p>
                </TooltipContent>
              )}
            </Tooltip>
          </TooltipProvider>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px] md:max-w-[650px] max-h-[85vh] overflow-y-auto bg-background border border-border">
        <DialogHeader>
          <DialogTitle>Create New Task</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              placeholder="Enter task title..."
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              autoFocus
              className="w-full"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              placeholder="Enter task description..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full resize-none"
            />
          </div>

          <div className="space-y-2">
            <Label>Task Color</Label>
            <div className="w-full">
              <ColorPicker
                value={color}
                onChange={setColor}
                size="md"
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

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Priority</Label>
              <Select value={priority} onValueChange={setPriority}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">
                    <div className="flex items-center gap-2">
                      <Flag className="h-3 w-3 text-green-500" />
                      Low
                    </div>
                  </SelectItem>
                  <SelectItem value="medium">
                    <div className="flex items-center gap-2">
                      <Flag className="h-3 w-3 text-yellow-500" />
                      Medium
                    </div>
                  </SelectItem>
                  <SelectItem value="high">
                    <div className="flex items-center gap-2">
                      <Flag className="h-3 w-3 text-red-500" />
                      High
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="due-date">Due Date</Label>
              <div className="relative">
                <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="due-date"
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  className="pl-10 w-full"
                />
              </div>
            </div>
          </div>

          <div className="text-sm text-muted-foreground">
            Adding to: <span className="font-medium">{listName}</span>
          </div>

          <div className="flex gap-2 pt-4">
            <Button type="submit" className="flex-1" disabled={isLoading}>
              {isLoading ? 'Creating...' : 'Create Task'}
            </Button>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => setOpen(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}