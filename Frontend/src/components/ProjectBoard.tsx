import { useState } from 'react';
import {
  DndContext,
  DragEndEvent,
  DragOverEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
  closestCenter,
} from '@dnd-kit/core';
import { 
  arrayMove, 
  SortableContext, 
  horizontalListSortingStrategy 
} from '@dnd-kit/sortable';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Plus, MoreHorizontal, Edit, Palette, FileText, Database, GitBranch, Bot, Users } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { ListColumn } from './ListColumn';
import { TaskCard } from './TaskCard';
import { CreateListDialog } from './CreateListDialog';
import { EditProjectDialog } from './EditProjectDialog';
import { ProjectMembersDialog } from './ProjectMembersDialog';
import { ColorIndicator } from '@/components/ui/color-picker';
import { ProjectWithLists, Task, ListWithTasks, Project } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';

interface ProjectBoardProps {
  project: ProjectWithLists;
  onTaskMove?: (taskId: string, newListUid: string) => void;
  onTaskClick?: (taskId: string) => void;
  onAddTask?: (listUid: string, taskData: any) => void;
  onAddList?: () => void;
  onDeleteList?: (listUid: string) => void;
  onDeleteTask?: (taskUid: string) => void;
  onProjectUpdate?: (updatedProject: Project) => void;
  onListUpdate?: (listUid: string, updatedList: ListWithTasks) => void;
  onTaskUpdate?: (taskUid: string, updatedTask: Task) => void;
  onListMove?: (listId: string, newPosition: number) => void;
}

export function ProjectBoard({
  project,
  onTaskMove,
  onTaskClick,
  onAddTask,
  onAddList,
  onDeleteList,
  onDeleteTask,
  onProjectUpdate,
  onListUpdate,
  onTaskUpdate,
  onListMove,
}: ProjectBoardProps) {
  const navigate = useNavigate();
  const [activeTask, setActiveTask] = useState<Task | null>(null);
  const [activeList, setActiveList] = useState<ListWithTasks | null>(null);
  const [showEditProjectDialog, setShowEditProjectDialog] = useState(false);
  const { toast } = useToast();
  
  // Use project.lists directly instead of local state, with fallback to empty array
  const lists = project.lists || [];

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event;
    
    // Check if we're dragging a task
    const task = lists
      .flatMap(list => list.tasks || [])
      .find(task => task.task_uid === active.id);
    
    if (task) {
      setActiveTask(task);
      return;
    }

    // Check if we're dragging a list
    const list = lists.find(list => list.list_uid === active.id);
    if (list) {
      setActiveList(list);
    }
  };

  const handleDragOver = (event: DragOverEvent) => {
    // For visual feedback during drag, we don't need to update state
    // The actual move will be handled in handleDragEnd
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    
    setActiveTask(null);
    setActiveList(null);
    
    if (!over) return;

    const activeId = active.id;
    const overId = over.id;

    if (activeId === overId) return;

    // Check if we're moving a list
    const isMovingList = lists.some(list => list.list_uid === activeId);
    
    if (isMovingList) {
      const oldIndex = lists.findIndex(list => list.list_uid === activeId);
      let newIndex = lists.findIndex(list => list.list_uid === overId);
      
      // If we didn't find the exact overId in lists, it might be a task within a list
      // In that case, find the list that contains this task
      if (newIndex === -1) {
        newIndex = lists.findIndex(list => 
          (list.tasks || []).some(task => task.task_uid === overId)
        );
      }
      
      if (oldIndex !== -1 && newIndex !== -1 && oldIndex !== newIndex) {
        try {
          // Convert array index to a proper position (0-based)
          const newPosition = newIndex;
          
          console.log('Moving list:', { 
            activeId, 
            overId, 
            oldIndex, 
            newIndex, 
            newPosition,
            listsCount: lists.length 
          });
          
          // Update position on the backend
          await apiService.updateListPosition(activeId as string, newPosition);
          
          // Call parent callback if provided
          if (onListMove) {
            onListMove(activeId as string, newPosition);
          }

          toast({
            title: 'List moved',
            description: 'List position updated successfully',
          });
        } catch (error) {
          console.error('Failed to update list position:', error);
          const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred';
          console.log('Full error details:', error);
          toast({
            title: 'Failed to move list',
            description: errorMessage,
            variant: 'destructive',
          });
        }
      } else {
        console.log('Skipping list move:', { oldIndex, newIndex, activeId, overId });
      }
      return;
    }

    // Handle task movement (existing logic)
    const activeTask = lists
      .flatMap(list => list.tasks || [])
      .find(task => task.task_uid === activeId);
    
    if (!activeTask) return;

    const overList = lists.find(list => list.list_uid === overId) ||
                    lists.find(list => 
                      (list.tasks || []).some(task => task.task_uid === overId)
                    );

    if (!overList) return;

    // Call the parent callback for task movement
    onTaskMove?.(activeTask.task_uid, overList.list_uid);
  };

  const handleEditProject = () => {
    setShowEditProjectDialog(true);
  };

  const handleCanvasClick = () => {
    navigate(`/project/${project.project_uid}/canvas`);
  };

  const handleNotesClick = () => {
    navigate(`/project/${project.project_uid}/notes`);
  };

  const handleDevAIClick = () => {
    navigate(`/project/${project.project_uid}/dev-ai`);
  };

  const handleProjectUpdated = (updatedProject: Project) => {
    if (onProjectUpdate) {
      onProjectUpdate(updatedProject);
    }
    setShowEditProjectDialog(false);
  };

  const handleListUpdate = (listUid: string, updatedList: ListWithTasks) => {
    // Call parent callback to update state
    if (onListUpdate) {
      onListUpdate(listUid, updatedList);
    }
  };

  const handleTaskUpdate = (taskUid: string, updatedTask: Task) => {
    // Call parent callback to update state
    if (onTaskUpdate) {
      onTaskUpdate(taskUid, updatedTask);
    }
  };

  return (
    <div className="flex-1 overflow-hidden flex flex-col">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={handleDragStart}
        onDragOver={handleDragOver}
        onDragEnd={handleDragEnd}
      >
        {/* Sticky Header */}
        <div className="flex-shrink-0 p-3 sm:p-6 pb-0 bg-background border-b">
          <div className="flex flex-col gap-3 sm:gap-4 mb-4 sm:mb-6">
            {/* Title Row */}
            <div className="flex items-center justify-between gap-2">
              <div className="flex items-center gap-2 sm:gap-3 flex-1 min-w-0">
                <ColorIndicator color={project.color || '#FFFFFF'} size="md" className="flex-shrink-0" />
                <h1 className="text-xl sm:text-2xl font-bold text-foreground truncate">{project.name}</h1>
              </div>
              <div className="flex items-center gap-2 flex-shrink-0">
                {/* Desktop: Members Button */}
                <ProjectMembersDialog 
                  projectId={project.project_uid}
                  trigger={
                    <Button
                      variant="outline"
                      size="sm"
                      className="hidden sm:flex items-center gap-2"
                    >
                      <Users className="h-4 w-4" />
                      <span className="hidden md:inline">Members</span>
                    </Button>
                  }
                />
                {/* Mobile: Members Icon */}
                <ProjectMembersDialog 
                  projectId={project.project_uid}
                  trigger={
                    <Button
                      variant="outline"
                      size="sm"
                      className="sm:hidden p-2"
                    >
                      <Users className="h-4 w-4" />
                    </Button>
                  }
                />
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 w-8 p-0 hover:bg-secondary"
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      className="cursor-pointer"
                      onClick={handleEditProject}
                    >
                      <Edit className="mr-2 h-4 w-4" />
                      Edit Project
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
            
            {/* Action Buttons Row - Scrollable on Mobile */}
            <div className="flex items-center gap-2 overflow-x-auto pb-2 sm:pb-0 -mx-3 px-3 sm:mx-0 sm:px-0">
              <Button
                variant="outline"
                size="sm"
                onClick={handleCanvasClick}
                className="flex items-center gap-2 whitespace-nowrap flex-shrink-0"
              >
                <Palette className="h-4 w-4" />
                <span className="hidden sm:inline">Canvas</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleNotesClick}
                className="flex items-center gap-2 whitespace-nowrap flex-shrink-0"
              >
                <FileText className="h-4 w-4" />
                <span className="hidden sm:inline">Notes</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleDevAIClick}
                className="flex items-center gap-2 whitespace-nowrap flex-shrink-0"
              >
                <Bot className="h-4 w-4" />
                <span className="hidden sm:inline">DEV AI</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => navigate(`/db-viewer?projectId=${project.project_uid}`)}
                className="hidden lg:flex items-center gap-2 whitespace-nowrap flex-shrink-0"
              >
                <Database className="h-4 w-4" />
                DB Viewer
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => navigate(`/flowchart?projectId=${project.project_uid}`)}
                className="hidden lg:flex items-center gap-2 whitespace-nowrap flex-shrink-0"
              >
                <GitBranch className="h-4 w-4" />
                Flowchart
              </Button>
            </div>
            
            {/* Description */}
            {project.description && (
              <p className="text-sm sm:text-base text-muted-foreground">{project.description}</p>
            )}
          </div>
        </div>

        {/* Scrollable Lists Section */}
        <div className="flex-1 overflow-auto p-3 sm:p-6 pt-4 sm:pt-6">
          <SortableContext items={lists.map(list => list.list_uid)} strategy={horizontalListSortingStrategy}>
            <div className="flex gap-3 sm:gap-6 pb-6 min-h-full -mx-3 px-3 sm:mx-0 sm:px-0">
              {lists.map((list) => (
                <ListColumn
                  key={list.list_uid}
                  list={list}
                  projectId={project.project_uid}
                  onTaskClick={onTaskClick}
                  onAddTask={onAddTask}
                  onListDelete={onDeleteList}
                  onListUpdate={handleListUpdate}
                  onTaskDelete={onDeleteTask}
                  onTaskUpdate={handleTaskUpdate}
                  isDraggable={true}
                />
              ))}
              
              <div className="flex-shrink-0">
                <CreateListDialog
                  projectUid={project.project_uid}
                  projectName={project.name}
                  onListCreate={() => onAddList?.()}
                />
              </div>
            </div>
          </SortableContext>
        </div>

        <DragOverlay>
          {activeTask ? (
            <div className="rotate-6 opacity-95">
              <TaskCard task={activeTask} />
            </div>
          ) : activeList ? (
            <div className="opacity-95 transform rotate-2">
              <ListColumn
                list={activeList}
                projectId={project.project_uid}
                onTaskClick={onTaskClick}
                onAddTask={onAddTask}
                onListDelete={onDeleteList}
                onListUpdate={handleListUpdate}
                onTaskDelete={onDeleteTask}
                onTaskUpdate={handleTaskUpdate}
                isDragOverlay={true}
              />
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>

      <EditProjectDialog
        project={project}
        open={showEditProjectDialog}
        onOpenChange={setShowEditProjectDialog}
        onProjectUpdated={handleProjectUpdated}
      />
    </div>
  );
}