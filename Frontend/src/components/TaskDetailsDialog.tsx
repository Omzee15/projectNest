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
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Calendar, Flag, User, Clock, Edit, Users, MessageSquare, Tag, Send, Trash2, X, Plus, Link, Check } from 'lucide-react';
import { Task, TaskComment, TaskCategory } from '@/types';
import { ColorIndicator } from '@/components/ui/color-picker';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { formatDistanceToNow } from 'date-fns';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';

const CATEGORY_COLORS = [
  '#EF4444', '#F97316', '#F59E0B', '#10B981',
  '#3B82F6', '#8B5CF6', '#EC4899', '#6B7280'
];

interface TaskDetailsDialogProps {
  task?: Task;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onEditTask: (task: Task) => void;
  projectUid?: string;
}

export function TaskDetailsDialog({ task, open, onOpenChange, onEditTask, projectUid }: TaskDetailsDialogProps) {
  const [comments, setComments] = useState<TaskComment[]>([]);
  const [categories, setCategories] = useState<TaskCategory[]>([]);
  const [projectCategories, setProjectCategories] = useState<TaskCategory[]>([]);
  const [newComment, setNewComment] = useState('');
  const [loadingComments, setLoadingComments] = useState(false);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [categoryPopoverOpen, setCategoryPopoverOpen] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState('');
  const [newCategoryColor, setNewCategoryColor] = useState('#3B82F6');
  const [linkCopied, setLinkCopied] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    if (open && task) {
      loadComments();
      loadCategories();
      if (projectUid) {
        loadProjectCategories();
      }
    }
  }, [open, task?.task_uid]);

  const loadComments = async () => {
    if (!task) return;
    try {
      setLoadingComments(true);
      const response = await apiService.getTaskComments(task.task_uid);
      setComments(response.data || []);
    } catch (error: any) {
      console.error('Failed to load comments:', error);
    } finally {
      setLoadingComments(false);
    }
  };

  const loadCategories = async () => {
    if (!task) return;
    try {
      setLoadingCategories(true);
      const response = await apiService.getTaskCategories(task.task_uid);
      setCategories(response.data || []);
    } catch (error: any) {
      console.error('Failed to load categories:', error);
    } finally {
      setLoadingCategories(false);
    }
  };

  const loadProjectCategories = async () => {
    if (!projectUid) return;
    try {
      const response = await apiService.getProjectCategories(projectUid);
      setProjectCategories(response.data || []);
    } catch (error: any) {
      console.error('Failed to load project categories:', error);
    }
  };

  const fetchProjectCategories = loadProjectCategories;

  const handleCopyLink = async () => {
    if (!task || !projectUid) return;
    
    const url = `${window.location.origin}/project/${projectUid}?task=${task.task_uid}`;
    
    try {
      await navigator.clipboard.writeText(url);
      setLinkCopied(true);
      toast({
        title: 'Link copied!',
        description: 'Shareable task link copied to clipboard',
      });
      
      // Reset the copied state after 2 seconds
      setTimeout(() => setLinkCopied(false), 2000);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to copy link to clipboard',
        variant: 'destructive',
      });
    }
  };

  const handleAddComment = async () => {
    if (!task || !newComment.trim()) return;

    try {
      await apiService.createTaskComment({
        task_uid: task.task_uid,
        content: newComment.trim(),
      });

      setNewComment('');
      loadComments();
      
      toast({
        title: 'Success',
        description: 'Comment added successfully',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to add comment',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteComment = async (commentUid: string) => {
    if (!confirm('Delete this comment?')) return;

    try {
      await apiService.deleteTaskComment(commentUid);
      loadComments();
      
      toast({
        title: 'Success',
        description: 'Comment deleted',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to delete comment',
        variant: 'destructive',
      });
    }
  };

  const handleToggleCategory = async (categoryUid: string) => {
    if (!task) return;

    const isAssigned = categories.some(c => c.category_uid === categoryUid);

    try {
      if (isAssigned) {
        await apiService.removeCategoryFromTask(task.task_uid, categoryUid);
      } else {
        // Assign this category (note: backend replaces all, so we need to send all UIDs)
        const newCategoryUids = [...categories.map(c => c.category_uid), categoryUid];
        await apiService.assignCategoriesToTask(task.task_uid, newCategoryUids);
      }

      loadCategories();
      
      toast({
        title: 'Success',
        description: isAssigned ? 'Category removed' : 'Category added',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to update categories',
        variant: 'destructive',
      });
    }
  };

  const handleCreateCategory = async () => {
    if (!projectUid || !newCategoryName.trim()) {
      toast({
        title: 'Error',
        description: 'Category name is required',
        variant: 'destructive',
      });
      return;
    }

    try {
      await apiService.createTaskCategory({
        project_uid: projectUid,
        name: newCategoryName.trim(),
        color: newCategoryColor,
      });

      setNewCategoryName('');
      setNewCategoryColor('#3B82F6');
      setShowCreateForm(false);
      loadProjectCategories();
      
      toast({
        title: 'Success',
        description: 'Category created successfully',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message || 'Failed to create category',
        variant: 'destructive',
      });
    }
  };

  if (!task) return null;

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getPriorityColor = (priority?: string) => {
    switch (priority?.toLowerCase()) {
      case 'high':
        return 'text-red-500';
      case 'medium':
        return 'text-yellow-500';
      case 'low':
        return 'text-green-500';
      default:
        return 'text-gray-500';
    }
  };

  const getStatusColor = (status: string, isCompleted: boolean) => {
    if (isCompleted) return 'border-green-500 text-green-700 bg-green-50';
    switch (status?.toLowerCase()) {
      case 'in-progress':
        return 'border-blue-500 text-blue-700 bg-blue-50';
      case 'todo':
        return 'border-gray-500 text-gray-700 bg-gray-50';
      case 'review':
        return 'border-purple-500 text-purple-700 bg-purple-50';
      default:
        return 'border-gray-500 text-gray-700 bg-gray-50';
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto overflow-x-hidden">
        <DialogHeader>
          <div className="flex flex-col gap-3 max-w-full">
            <div className="flex items-start justify-between gap-4 pr-8 max-w-full">
              <div className="flex items-center gap-3 flex-1 min-w-0 overflow-hidden">
                <ColorIndicator color={task.color || '#FFFFFF'} size="md" className="flex-shrink-0" />
                <DialogTitle className="text-lg sm:text-xl font-semibold leading-tight">
                  {task.title}
                </DialogTitle>
              </div>
            </div>
            <div className="flex items-center justify-between gap-4">
              <DialogDescription className="text-sm sm:text-base">
                Task Details and Information
              </DialogDescription>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleCopyLink}
                  className="flex items-center gap-2 flex-shrink-0"
                  title="Copy shareable link"
                >
                  {linkCopied ? (
                    <Check className="h-4 w-4 text-green-500" />
                  ) : (
                    <Link className="h-4 w-4" />
                  )}
                  <span className="hidden sm:inline">{linkCopied ? 'Copied!' : 'Share'}</span>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onEditTask(task)}
                  className="flex items-center gap-2 flex-shrink-0"
                >
                  <Edit className="h-4 w-4" />
                  <span className="hidden sm:inline">Edit</span>
                </Button>
              </div>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-6 max-w-full overflow-hidden">
          {/* Comments Display Section - Moved to top */}
          {comments.length > 0 && (
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <MessageSquare className="h-4 w-4" />
                Comments ({comments.length})
              </h3>
              <div className="space-y-3 max-h-48 overflow-y-auto">
                {comments.map((comment) => (
                  <div
                    key={comment.comment_uid}
                    className="bg-muted/50 rounded-lg p-3 space-y-2"
                  >
                    <div className="flex items-start justify-between gap-2">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-medium">
                            {comment.user_name || `User ${comment.user_id}`}
                          </span>
                          <span className="text-xs text-muted-foreground">
                            {formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}
                          </span>
                        </div>
                        <p className="text-sm mt-1 whitespace-pre-wrap break-words">
                          {comment.content}
                        </p>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteComment(comment.comment_uid)}
                        className="flex-shrink-0"
                      >
                        <Trash2 className="h-3 w-3 text-destructive" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Status and Priority Section */}
          <div className="flex items-center gap-4 flex-wrap max-w-full">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-muted-foreground">Status:</span>
              <Badge 
                variant="outline" 
                className={`text-sm px-3 py-1 ${getStatusColor(task.status, task.is_completed)}`}
              >
                {task.is_completed ? 'Completed' : task.status || 'To Do'}
              </Badge>
            </div>
            
            {task.priority && (
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium text-muted-foreground">Priority:</span>
                <Badge variant="secondary" className="text-sm px-3 py-1">
                  <Flag className={`w-3 h-3 mr-1 ${getPriorityColor(task.priority)}`} />
                  {task.priority}
                </Badge>
              </div>
            )}
          </div>

          {/* Description */}
          {task.description && (
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-muted-foreground">Description</h3>
              <div className="bg-muted/50 p-4 rounded-lg overflow-hidden">
                <p className="text-sm leading-relaxed whitespace-pre-wrap break-words overflow-wrap-anywhere">
                  {task.description}
                </p>
              </div>
            </div>
          )}

          {/* Assignees */}
          {task.assignees && task.assignees.length > 0 && (
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-muted-foreground">Assigned To</h3>
              <div className="flex flex-wrap gap-2">
                {task.assignees.map((assignee) => (
                  <Badge
                    key={assignee.user_uid}
                    variant="secondary"
                    className="flex items-center gap-2 py-1.5 px-3"
                  >
                    <Avatar className="h-5 w-5 border border-background">
                      <AvatarFallback className="text-[10px] bg-primary/10 text-primary">
                        {assignee.name.split(' ').map((n) => n[0]).join('').toUpperCase().slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                    <span className="text-xs font-medium">{assignee.name}</span>
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {/* Categories */}
          <div className="space-y-3 border-t pt-4">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Tag className="h-4 w-4" />
                Categories
              </h3>
              <Popover open={categoryPopoverOpen} onOpenChange={setCategoryPopoverOpen}>
                <PopoverTrigger asChild>
                  <Button variant="ghost" size="sm">
                    <Plus className="h-4 w-4 mr-1" />
                    Add
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-80" align="end">
                  <div className="space-y-3">
                    {/* Create Form */}
                    {showCreateForm ? (
                      <div className="space-y-3">
                        <div className="flex items-center justify-between">
                          <h4 className="font-medium text-sm">Create Category</h4>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              setShowCreateForm(false);
                              setNewCategoryName('');
                            }}
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                        
                        <div className="space-y-2">
                          <Label htmlFor="category-name">Name</Label>
                          <Input
                            id="category-name"
                            value={newCategoryName}
                            onChange={(e) => setNewCategoryName(e.target.value)}
                            placeholder="e.g., Bug, Feature, Priority"
                            autoFocus
                          />
                        </div>
                        
                        <div className="space-y-2">
                          <Label>Color</Label>
                          <div className="flex flex-wrap gap-2">
                            {CATEGORY_COLORS.map((colorOption) => (
                              <button
                                key={colorOption}
                                type="button"
                                className={`w-8 h-8 rounded-full transition-all hover:scale-110 ${
                                  newCategoryColor === colorOption ? 'ring-2 ring-offset-2 ring-primary' : ''
                                }`}
                                style={{ backgroundColor: colorOption }}
                                onClick={() => setNewCategoryColor(colorOption)}
                              />
                            ))}
                          </div>
                        </div>
                        
                        <div className="flex gap-2">
                          <Button
                            onClick={handleCreateCategory}
                            className="flex-1"
                            disabled={!newCategoryName.trim()}
                          >
                            Create
                          </Button>
                          <Button
                            variant="outline"
                            onClick={() => {
                              setShowCreateForm(false);
                              setNewCategoryName('');
                            }}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="flex items-center justify-between">
                          <h4 className="font-medium text-sm">Manage Categories</h4>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setShowCreateForm(true)}
                            className="h-7"
                          >
                            <Plus className="h-3 w-3 mr-1" />
                            Create New
                          </Button>
                        </div>
                        
                        {projectCategories.length === 0 ? (
                          <p className="text-sm text-muted-foreground text-center py-4">
                            No categories yet. Create one to get started!
                          </p>
                        ) : (
                          <div className="space-y-2">
                            <p className="text-xs text-muted-foreground">
                              Click to add or remove:
                            </p>
                            <div className="flex flex-wrap gap-2 max-h-60 overflow-y-auto">
                              {projectCategories.map((category) => {
                                const isAssigned = categories.some(c => c.category_uid === category.category_uid);
                                return (
                                  <Badge
                                    key={category.category_uid}
                                    variant={isAssigned ? "default" : "outline"}
                                    className="cursor-pointer transition-all hover:scale-105"
                                    style={isAssigned ? { backgroundColor: category.color, borderColor: category.color, color: 'white' } : { borderColor: category.color, color: category.color }}
                                    onClick={() => handleToggleCategory(category.category_uid)}
                                  >
                                    {category.name}
                                    {isAssigned && <X className="ml-1 h-3 w-3" />}
                                  </Badge>
                                );
                              })}
                            </div>
                          </div>
                        )}
                      </>
                    )}
                  </div>
                </PopoverContent>
              </Popover>
            </div>
            
            {/* Current Categories */}
            <div className="flex flex-wrap gap-2">
              {loadingCategories ? (
                <span className="text-sm text-muted-foreground">Loading...</span>
              ) : categories.length === 0 ? (
                <span className="text-sm text-muted-foreground italic">No categories assigned</span>
              ) : (
                categories.map((category) => (
                  <Badge
                    key={category.category_uid}
                    style={{ backgroundColor: category.color, borderColor: category.color }}
                    className="text-white"
                  >
                    {category.name}
                  </Badge>
                ))
              )}
            </div>
          </div>

          {/* Due Date */}
          {task.due_date && (
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-muted-foreground">Due Date</h3>
              <div className="flex items-center gap-2 text-sm">
                <Calendar className="h-4 w-4 text-muted-foreground" />
                <span>{formatDate(task.due_date)}</span>
                {new Date(task.due_date) < new Date() && !task.is_completed && (
                  <Badge variant="destructive" className="text-xs">
                    Overdue
                  </Badge>
                )}
              </div>
            </div>
          )}

          {/* Completion Information */}
          {task.is_completed && task.completed_at && (
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-muted-foreground">Completed</h3>
              <div className="flex items-center gap-2 text-sm">
                <Clock className="h-4 w-4 text-green-500" />
                <span>{formatDateTime(task.completed_at)}</span>
                <span className="text-muted-foreground">
                  ({formatDistanceToNow(new Date(task.completed_at), { addSuffix: true })})
                </span>
              </div>
            </div>
          )}

          {/* Created By and Timestamps */}
          <div className="space-y-3 border-t pt-4">
            <h3 className="text-sm font-medium text-muted-foreground">Task Information</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Created by:</span>
                  <span className="text-muted-foreground">{task.created_by}</span>
                </div>
                
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Created:</span>
                  <span className="text-muted-foreground">{formatDateTime(task.created_at)}</span>
                </div>
              </div>

              {task.updated_at && task.updated_by && (
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Updated by:</span>
                    <span className="text-muted-foreground">{task.updated_by}</span>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Updated:</span>
                    <span className="text-muted-foreground">{formatDateTime(task.updated_at)}</span>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Add Comment Form - At the bottom */}
          <div className="space-y-3 border-t pt-4">
            <h3 className="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <MessageSquare className="h-4 w-4" />
              Add Comment
            </h3>
            <div className="flex gap-2">
              <Textarea
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="Add a comment..."
                rows={2}
                className="flex-1"
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
                    handleAddComment();
                  }
                }}
              />
              <Button
                onClick={handleAddComment}
                disabled={!newComment.trim()}
                size="sm"
              >
                <Send className="h-4 w-4" />
              </Button>
            </div>
            {loadingComments && (
              <div className="text-center py-2 text-sm text-muted-foreground">
                Loading comments...
              </div>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}