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
import { Calendar, Flag, User, Clock, Edit } from 'lucide-react';
import { Task } from '@/types';
import { ColorIndicator } from '@/components/ui/color-picker';
import { formatDistanceToNow } from 'date-fns';

interface TaskDetailsDialogProps {
  task?: Task;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onEditTask: (task: Task) => void;
}

export function TaskDetailsDialog({ task, open, onOpenChange, onEditTask }: TaskDetailsDialogProps) {
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
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <div className="flex flex-col gap-3">
            <div className="flex items-start justify-between gap-4 pr-8">
              <div className="flex items-center gap-3 flex-1 min-w-0">
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
        </DialogHeader>

        <div className="space-y-6">
          {/* Status and Priority Section */}
          <div className="flex items-center gap-4 flex-wrap">
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
              <div className="bg-muted/50 p-4 rounded-lg">
                <p className="text-sm leading-relaxed whitespace-pre-wrap">
                  {task.description}
                </p>
              </div>
            </div>
          )}

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