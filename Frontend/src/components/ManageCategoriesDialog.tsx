import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Plus, Trash2, Edit2, X } from 'lucide-react';
import { TaskCategory, TaskCategoryRequest, TaskCategoryUpdateRequest } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';

interface ManageCategoriesDialogProps {
  projectUid: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCategoriesUpdated?: () => void;
}

const CATEGORY_COLORS = [
  '#EF4444', '#F97316', '#F59E0B', '#10B981', 
  '#3B82F6', '#8B5CF6', '#EC4899', '#6B7280'
];

export function ManageCategoriesDialog({ 
  projectUid, 
  open, 
  onOpenChange,
  onCategoriesUpdated 
}: ManageCategoriesDialogProps) {
  const [categories, setCategories] = useState<TaskCategory[]>([]);
  const [loading, setLoading] = useState(false);
  const [editingCategory, setEditingCategory] = useState<TaskCategory | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  
  // Form state
  const [name, setName] = useState('');
  const [color, setColor] = useState(CATEGORY_COLORS[0]);
  const [description, setDescription] = useState('');
  
  const { toast } = useToast();

  useEffect(() => {
    if (open) {
      loadCategories();
    }
  }, [open, projectUid]);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const response = await apiService.getProjectCategories(projectUid);
      setCategories(response.data || []);
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to load categories',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const resetForm = () => {
    setName('');
    setColor(CATEGORY_COLORS[0]);
    setDescription('');
    setEditingCategory(null);
    setShowCreateForm(false);
  };

  const handleCreate = async () => {
    if (!name.trim()) {
      toast({
        title: 'Error',
        description: 'Category name is required',
        variant: 'destructive',
      });
      return;
    }

    try {
      const request: TaskCategoryRequest = {
        project_uid: projectUid,
        name: name.trim(),
        color,
        description: description.trim() || undefined,
      };

      await apiService.createTaskCategory(request);
      
      toast({
        title: 'Success',
        description: 'Category created successfully',
      });

      resetForm();
      loadCategories();
      onCategoriesUpdated?.();
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to create category',
        variant: 'destructive',
      });
    }
  };

  const handleUpdate = async () => {
    if (!editingCategory || !name.trim()) return;

    try {
      const request: TaskCategoryUpdateRequest = {
        name: name.trim(),
        color,
        description: description.trim() || undefined,
      };

      await apiService.updateTaskCategory(editingCategory.category_uid, request);
      
      toast({
        title: 'Success',
        description: 'Category updated successfully',
      });

      resetForm();
      loadCategories();
      onCategoriesUpdated?.();
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to update category',
        variant: 'destructive',
      });
    }
  };

  const handleDelete = async (categoryUid: string) => {
    if (!confirm('Are you sure you want to delete this category? It will be removed from all tasks.')) {
      return;
    }

    try {
      await apiService.deleteTaskCategory(categoryUid);
      
      toast({
        title: 'Success',
        description: 'Category deleted successfully',
      });

      loadCategories();
      onCategoriesUpdated?.();
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to delete category',
        variant: 'destructive',
      });
    }
  };

  const startEdit = (category: TaskCategory) => {
    setEditingCategory(category);
    setName(category.name);
    setColor(category.color);
    setDescription(category.description || '');
    setShowCreateForm(true);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Manage Categories</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Create/Edit Form */}
          {showCreateForm ? (
            <div className="border rounded-lg p-4 space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="font-semibold">
                  {editingCategory ? 'Edit Category' : 'New Category'}
                </h3>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={resetForm}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>

              <div className="space-y-3">
                <div>
                  <Label htmlFor="name">Name *</Label>
                  <Input
                    id="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="e.g., Bug, Feature, Documentation"
                    maxLength={100}
                  />
                </div>

                <div>
                  <Label>Color</Label>
                  <div className="flex gap-2 mt-2">
                    {CATEGORY_COLORS.map((c) => (
                      <button
                        key={c}
                        type="button"
                        className={`w-8 h-8 rounded-full border-2 transition-all ${
                          color === c ? 'border-primary scale-110' : 'border-transparent'
                        }`}
                        style={{ backgroundColor: c }}
                        onClick={() => setColor(c)}
                      />
                    ))}
                  </div>
                </div>

                <div>
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="Optional description"
                    rows={2}
                  />
                </div>

                <div className="flex gap-2">
                  <Button
                    onClick={editingCategory ? handleUpdate : handleCreate}
                    className="flex-1"
                  >
                    {editingCategory ? 'Update' : 'Create'}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={resetForm}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          ) : (
            <Button
              onClick={() => setShowCreateForm(true)}
              className="w-full"
              variant="outline"
            >
              <Plus className="h-4 w-4 mr-2" />
              Create New Category
            </Button>
          )}

          {/* Categories List */}
          <div className="space-y-2">
            <h3 className="text-sm font-medium text-muted-foreground">
              Existing Categories ({categories.length})
            </h3>
            
            {loading ? (
              <div className="text-center py-4 text-muted-foreground">Loading...</div>
            ) : categories.length === 0 ? (
              <div className="text-center py-4 text-muted-foreground">
                No categories yet. Create one to get started!
              </div>
            ) : (
              <div className="space-y-2">
                {categories.map((category) => (
                  <div
                    key={category.category_uid}
                    className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <div
                        className="w-4 h-4 rounded-full flex-shrink-0"
                        style={{ backgroundColor: category.color }}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="font-medium">{category.name}</div>
                        {category.description && (
                          <div className="text-sm text-muted-foreground truncate">
                            {category.description}
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="flex gap-2 flex-shrink-0">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => startEdit(category)}
                      >
                        <Edit2 className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(category.category_uid)}
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
