import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ColorPicker } from '@/components/ui/color-picker';
import { Plus, Lock } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { apiService } from '@/services/api';
import { COLORS } from '@/types';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

interface CreateListDialogProps {
  projectUid: string;
  projectName: string;
  onListCreate?: () => void;
  trigger?: React.ReactNode;
  canWrite?: boolean;
}

export function CreateListDialog({ projectUid, projectName, onListCreate, trigger, canWrite = true }: CreateListDialogProps) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [color, setColor] = useState<string>(COLORS.WHITE);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      toast({
        title: 'Error',
        description: 'List name is required',
        variant: 'destructive',
      });
      return;
    }

    setIsLoading(true);

    try {
      const listData = {
        project_uid: projectUid,
        name: name.trim(),
        color,
        position: 0, // Let backend auto-calculate position
      };

      console.log('[CreateListDialog] Creating list with data:', listData);

      await apiService.createList(listData);
      
      // Reset form
      setName('');
      setColor(COLORS.WHITE);
      setOpen(false);
      setIsLoading(false);

      toast({
        title: 'List created',
        description: `List "${name}" added to ${projectName}`,
      });

      // Notify parent to refresh data
      onListCreate?.();

    } catch (error) {
      console.error('[CreateListDialog] Failed to create list:', error);
      setIsLoading(false);
      toast({
        title: 'Failed to create list',
        description: error instanceof Error ? error.message : 'An unexpected error occurred',
        variant: 'destructive',
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={(newOpen) => {
      if (!canWrite && newOpen) {
        toast({
          title: 'Access Denied',
          description: 'You do not have permission to create lists. Only users with edit access can create lists.',
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
                  className="w-72 h-12 border-2 border-dashed border-border hover:border-border-light hover:bg-muted/50"
                  disabled={!canWrite}
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add another list
                </Button>
              </TooltipTrigger>
              {!canWrite && (
                <TooltipContent>
                  <p>You have view-only access. Only users with edit permission can create lists.</p>
                </TooltipContent>
              )}
            </Tooltip>
          </TooltipProvider>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-md bg-background border border-border">
        <DialogHeader>
          <DialogTitle>Create New List</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">List Name *</Label>
            <Input
              id="name"
              placeholder="Enter list name..."
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <Label>List Color</Label>
            <ColorPicker
              value={color}
              onChange={setColor}
              size="md"
            />
          </div>

          <div className="text-sm text-muted-foreground">
            Adding to: <span className="font-medium">{projectName}</span>
          </div>

          <div className="flex gap-2 pt-4">
            <Button type="submit" className="flex-1">
              Create List
            </Button>
            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}