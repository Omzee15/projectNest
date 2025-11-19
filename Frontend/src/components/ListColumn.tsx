import { useDroppable } from '@dnd-kit/core';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Card, CardHeader, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
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
import { Plus, MoreHorizontal, Trash2, Loader2, Edit, GripVertical } from 'lucide-react';
import { TaskCard } from './TaskCard';
import { CreateTaskDialog } from './CreateTaskDialog';
import { EditListDialog } from './EditListDialog';
import { EditTaskDialog } from './EditTaskDialog';
import { ListWithTasks, Task } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { ColorIndicator } from '@/components/ui/color-picker';
import { useState } from 'react';

interface ListColumnProps {
  list: ListWithTasks;
  projectId: string;
  onTaskClick?: (taskId: string) => void;
  onAddTask?: (listUid: string, taskData: any) => void;
  onListDelete?: (listUid: string) => void;
  onListUpdate?: (listUid: string, updatedList: ListWithTasks) => void;
  onTaskDelete?: (taskUid: string) => void;
  onTaskUpdate?: (taskUid: string, updatedTask: Task) => void;
  isDraggable?: boolean;
  isDragOverlay?: boolean;
}

export function ListColumn({ 
  list,
  projectId,
  onTaskClick, 
  onAddTask, 
  onListDelete, 
  onListUpdate, 
  onTaskDelete, 
  onTaskUpdate,
  isDraggable = false,
  isDragOverlay = false
}: ListColumnProps) {
  const { setNodeRef: setDroppableRef } = useDroppable({
    id: list.list_uid,
  });
  
  const {
    attributes,
    listeners,
    setNodeRef: setSortableRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: list.list_uid,
    disabled: !isDraggable || isDragOverlay,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const [isDeleting, setIsDeleting] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showEditListDialog, setShowEditListDialog] = useState(false);
  const [showEditTaskDialog, setShowEditTaskDialog] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | undefined>(undefined);
  const { toast } = useToast();

  const taskIds = (list.tasks || []).map(task => task.task_uid);

  const handleDeleteList = async () => {
    try {
      setIsDeleting(true);
      await onListDelete?.(list.list_uid);
      setShowDeleteDialog(false);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete list. Please try again.',
        variant: 'destructive',
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const handleEditList = () => {
    setShowEditListDialog(true);
  };

  const handleListUpdated = (updatedList: any) => {
    if (onListUpdate) {
      // Preserve the tasks when updating list
      onListUpdate(list.list_uid, { ...updatedList, tasks: list.tasks || [] });
    }
  };

  const handleTaskEdit = (task: Task) => {
    setEditingTask(task);
    setShowEditTaskDialog(true);
  };

  const handleTaskUpdated = (updatedTask: Task) => {
    if (onTaskUpdate) {
      onTaskUpdate(updatedTask.task_uid, updatedTask);
    }
    setShowEditTaskDialog(false);
    setEditingTask(undefined);
  };

  const handleTaskUpdate = (taskUid: string, updatedTask: Task) => {
    if (onTaskUpdate) {
      onTaskUpdate(taskUid, updatedTask);
    }
  };

  // Combine the refs for both sortable and droppable functionality
  const combineRefs = (node: HTMLElement | null) => {
    setDroppableRef(node);
    setSortableRef(node);
  };

  return (
    <div 
      ref={combineRefs}
      style={style}
      {...(isDraggable ? attributes : {})}
      className={isDragOverlay ? "pointer-events-none" : ""}
    >
      <Card className="w-72 flex-shrink-0 bg-muted/30 border-border-light">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              {isDraggable && !isDragOverlay && (
                <button
                  {...listeners}
                  className="cursor-grab hover:cursor-grabbing p-1 hover:bg-secondary rounded"
                >
                  <GripVertical className="h-4 w-4 text-muted-foreground" />
                </button>
              )}
              <ColorIndicator color={list.color || '#FFFFFF'} size="sm" />
              <h3 className="font-semibold text-sm text-foreground">
                {list.name}
              </h3>
            </div>
          <div className="flex items-center gap-1">
            <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
              {(list.tasks || []).length}
            </span>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0 hover:bg-secondary"
                >
                  <MoreHorizontal className="h-3 w-3" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  className="cursor-pointer"
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    handleEditList();
                  }}
                >
                  <Edit className="mr-2 h-4 w-4" />
                  Edit List
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
                  Delete List
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            
            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete List</AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to delete "{list.name}"? This will also delete all tasks in this list. This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={handleDeleteList}
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    disabled={isDeleting}
                  >
                    {isDeleting ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Deleting...
                      </>
                    ) : (
                      'Delete List'
                    )}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </div>
      </CardHeader>

      <CardContent className="pt-0 pb-3 space-y-3">
        <div className="space-y-2 min-h-2">
          <SortableContext items={taskIds} strategy={verticalListSortingStrategy}>
            {(list.tasks || []).map((task) => (
              <TaskCard
                key={task.task_uid}
                task={task}
                onClick={() => onTaskClick?.(task.task_uid)}
                onTaskDelete={onTaskDelete}
                onTaskUpdate={handleTaskUpdate}
                onTaskEdit={handleTaskEdit}
              />
            ))}
          </SortableContext>
        </div>

        <CreateTaskDialog
          listUid={list.list_uid}
          listName={list.name}
          projectId={projectId}
          onTaskCreate={(taskData) => onAddTask?.(list.list_uid, taskData)}
        />
      </CardContent>

      <EditListDialog
        list={list}
        open={showEditListDialog}
        onOpenChange={setShowEditListDialog}
        onListUpdated={handleListUpdated}
      />

      <EditTaskDialog
        task={editingTask}
        projectId={projectId}
        open={showEditTaskDialog}
        onOpenChange={(open) => {
          setShowEditTaskDialog(open);
          if (!open) setEditingTask(undefined);
        }}
        onTaskUpdated={handleTaskUpdated}
      />
      </Card>
    </div>
  );
}