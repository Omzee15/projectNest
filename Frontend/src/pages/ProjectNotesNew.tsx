import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { ProjectNotesCollaborative } from '@/components/ProjectNotesCollaborative';
import { apiService } from '@/services/api';
import { toast } from '@/hooks/use-toast';
import { NoteRequest, NoteUpdateRequest } from '@/types';

const ProjectNotes: React.FC = () => {
  const { projectUid } = useParams<{ projectUid: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Fetch notes data
  const { data: notesData, isLoading } = useQuery({
    queryKey: ['notes', projectUid],
    queryFn: () => apiService.getNotesByProject(projectUid!),
    enabled: !!projectUid,
  });

  // Create note mutation
  const createNoteMutation = useMutation({
    mutationFn: (noteData: NoteRequest) => 
      apiService.createNote(projectUid!, noteData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      toast({
        title: 'Success',
        description: 'Note created successfully',
      });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to create note',
        variant: 'destructive',
      });
    },
  });

  // Update note mutation
  const updateNoteMutation = useMutation({
    mutationFn: ({ noteUid, noteData }: { noteUid: string; noteData: NoteUpdateRequest }) =>
      apiService.partialUpdateNote(noteUid, noteData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to save note',
        variant: 'destructive',
      });
    },
  });

  // Delete note mutation
  const deleteNoteMutation = useMutation({
    mutationFn: (noteUid: string) => apiService.deleteNote(noteUid),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      toast({
        title: 'Success',
        description: 'Note deleted successfully',
      });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to delete note',
        variant: 'destructive',
      });
    },
  });

  const handleCreateNote = (noteData: NoteRequest) => {
    createNoteMutation.mutate(noteData);
  };

  const handleUpdateNote = (noteUid: string, noteData: NoteUpdateRequest) => {
    updateNoteMutation.mutate({ noteUid, noteData });
  };

  const handleDeleteNote = (noteUid: string) => {
    deleteNoteMutation.mutate(noteUid);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-pulse text-lg">Loading notes...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate(`/project/${projectUid}`)}
              className="flex items-center gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Project
            </Button>
            <div className="border-l h-6"></div>
            <h1 className="text-xl font-semibold">Collaborative Notes</h1>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-6 py-8">
        <ProjectNotesCollaborative
          projectUid={projectUid!}
          notes={notesData?.data?.notes || []}
          onCreateNote={handleCreateNote}
          onUpdateNote={handleUpdateNote}
          onDeleteNote={handleDeleteNote}
        />
      </div>
    </div>
  );
};

export default ProjectNotes;