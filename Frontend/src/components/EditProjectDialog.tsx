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
import { Loader2 } from 'lucide-react';
import { Project, ProjectUpdateRequest } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';

interface EditProjectDialogProps {
  project?: Project;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onProjectUpdated: (updatedProject: Project) => void;
}

export function EditProjectDialog({ project, open, onOpenChange, onProjectUpdated }: EditProjectDialogProps) {
  const [formData, setFormData] = useState<ProjectUpdateRequest>({
    name: '',
    description: '',
    status: '',
    color: '#FFFFFF',
    start_date: '',
    end_date: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  // Reset form when project changes
  useEffect(() => {
    if (project && open) {
      setFormData({
        name: project.name,
        description: project.description || '',
        status: project.status,
        color: project.color || '#FFFFFF',
        start_date: project.start_date ? project.start_date.split('T')[0] : '',
        end_date: project.end_date ? project.end_date.split('T')[0] : '',
      });
    }
  }, [project, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!project) return;

    try {
      setIsLoading(true);

      // Filter out empty values and prepare the update request
      const updates: ProjectUpdateRequest = {};
      if (formData.name?.trim() && formData.name !== project.name) {
        updates.name = formData.name.trim();
      }
      if (formData.description !== project.description) {
        updates.description = formData.description || undefined;
      }
      if (formData.status && formData.status !== project.status) {
        updates.status = formData.status;
      }
      if (formData.color && formData.color !== project.color) {
        updates.color = formData.color;
      }
      if (formData.start_date !== (project.start_date ? project.start_date.split('T')[0] : '')) {
        updates.start_date = formData.start_date || undefined;
      }
      if (formData.end_date !== (project.end_date ? project.end_date.split('T')[0] : '')) {
        updates.end_date = formData.end_date || undefined;
      }

      // Only send request if there are changes
      if (Object.keys(updates).length === 0) {
        toast({
          title: 'No changes',
          description: 'No changes were made to the project.',
        });
        onOpenChange(false);
        return;
      }

      const response = await apiService.partialUpdateProject(project.project_uid, updates);
      
      onProjectUpdated(response.data);
      onOpenChange(false);
      
      toast({
        title: 'Project updated',
        description: 'Project has been successfully updated.',
      });
    } catch (error) {
      console.error('Failed to update project:', error);
      toast({
        title: 'Error',
        description: 'Failed to update project. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (field: keyof ProjectUpdateRequest, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value === '' ? undefined : value
    }));
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[700px] lg:max-w-[800px] max-h-[90vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <DialogTitle>Edit Project</DialogTitle>
          <DialogDescription>
            Update project details.
          </DialogDescription>
        </DialogHeader>
        
        <div className="flex-1 overflow-y-auto mt-4">
            <form onSubmit={handleSubmit} className="space-y-4 pr-2">
              <div className="space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  placeholder="Enter project name"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => handleInputChange('description', e.target.value)}
                  placeholder="Enter project description (optional)"
                  rows={3}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="status">Status</Label>
                  <Select
                    value={formData.status}
                    onValueChange={(value) => handleInputChange('status', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="active">Active</SelectItem>
                      <SelectItem value="inactive">Inactive</SelectItem>
                      <SelectItem value="completed">Completed</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label>Color</Label>
                  <ColorPicker
                    value={formData.color || '#FFFFFF'}
                    onChange={(color) => handleInputChange('color', color)}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="start_date">Start Date</Label>
                  <Input
                    id="start_date"
                    type="date"
                    value={formData.start_date}
                    onChange={(e) => handleInputChange('start_date', e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="end_date">End Date</Label>
                  <Input
                    id="end_date"
                    type="date"
                    value={formData.end_date}
                    onChange={(e) => handleInputChange('end_date', e.target.value)}
                  />
                </div>
              </div>

              <DialogFooter className="mt-6">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => onOpenChange(false)}
                  disabled={isLoading}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Updating...
                    </>
                  ) : (
                    'Update Project'
                  )}
                </Button>
              </DialogFooter>
            </form>
          </div>
      </DialogContent>
    </Dialog>
  );
}