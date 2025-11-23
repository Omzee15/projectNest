import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { 
  Plus, 
  Edit, 
  Trash2, 
  FileText,
  Calendar,
  MoreVertical,
  Users
} from 'lucide-react';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Note, NoteRequest, NoteUpdateRequest } from '@/types';
import { CollaborativeEditor } from './CollaborativeEditor';
import { CollaboratorAvatars, ConnectionStatus, CollaboratorPresence } from './CollaboratorPresence';
import { getCurrentUser } from '@/utils/auth';

interface ProjectNotesCollaborativeProps {
  projectUid: string;
  notes: Note[];
  onCreateNote?: (noteData: NoteRequest) => void;
  onUpdateNote?: (noteUid: string, noteData: NoteUpdateRequest) => void;
  onDeleteNote?: (noteUid: string) => void;
  readOnly?: boolean;
}

interface NoteFormData {
  title: string;
  content: string;
}

const getUserColor = (userId: string): string => {
  const colors = [
    '#3B82F6', // blue
    '#10B981', // green
    '#F59E0B', // yellow
    '#EF4444', // red
    '#8B5CF6', // purple
    '#EC4899', // pink
    '#14B8A6', // teal
  ];
  
  // Generate a consistent color based on user ID
  const hash = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
  return colors[hash % colors.length];
};

export const ProjectNotesCollaborative: React.FC<ProjectNotesCollaborativeProps> = ({
  projectUid,
  notes,
  onCreateNote,
  onUpdateNote,
  onDeleteNote,
  readOnly = false
}) => {
  const user = getCurrentUser();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [formData, setFormData] = useState<NoteFormData>({ title: '', content: '' });
  const [connectionStatus, setConnectionStatus] = useState<'connected' | 'connecting' | 'disconnected'>('connecting');
  const [collaborators, setCollaborators] = useState<Record<string, CollaboratorPresence[]>>({});

  const handleCreateNote = () => {
    if (!formData.title.trim() || !onCreateNote) return;

    onCreateNote({
      title: formData.title.trim(),
      content: { blocks: [] }, // Empty content for new notes
    });

    setFormData({ title: '', content: '' });
    setIsCreateDialogOpen(false);
  };

  const handleSaveNote = (noteUid: string) => (content: string, html: string) => {
    if (!onUpdateNote) return;

    try {
      const contentObj = JSON.parse(content);
      onUpdateNote(noteUid, {
        content: contentObj,
        content_html: html,
        version: notes.find(n => n.note_uid === noteUid)?.version,
        last_modified_by: user?.user_uid,
      });
    } catch (error) {
      console.error('Failed to save note:', error);
    }
  };

  const handleDeleteNote = (noteUid: string) => {
    if (!onDeleteNote) return;
    onDeleteNote(noteUid);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const resetForm = () => {
    setFormData({ title: '', content: '' });
    setEditingNote(null);
    setIsCreateDialogOpen(false);
  };

  const getDocumentId = (noteUid: string) => {
    return `project-${projectUid}-note-${noteUid}`;
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <CardTitle className="text-lg flex items-center gap-2">
              <FileText className="w-5 h-5" />
              Project Notes
              <Badge variant="secondary" className="ml-2">
                {notes.length}
              </Badge>
            </CardTitle>
            <ConnectionStatus status={connectionStatus} />
          </div>
          
          {!readOnly && (
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
              <DialogTrigger asChild>
                <Button size="sm">
                  <Plus className="w-4 h-4 mr-1" />
                  New Note
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create New Note</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                  <div>
                    <label className="text-sm font-medium">Title</label>
                    <Input
                      value={formData.title}
                      onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
                      placeholder="Enter note title..."
                      className="mt-1"
                    />
                  </div>
                  <div className="flex justify-end gap-2">
                    <Button variant="outline" onClick={resetForm}>
                      Cancel
                    </Button>
                    <Button 
                      onClick={handleCreateNote}
                      disabled={!formData.title.trim()}
                    >
                      Create Note
                    </Button>
                  </div>
                </div>
              </DialogContent>
            </Dialog>
          )}
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="space-y-4">
          {notes.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <FileText className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p className="text-lg font-medium">No notes yet</p>
              <p className="text-sm">Create your first note to get started with collaborative editing.</p>
            </div>
          ) : (
            notes.map((note) => (
              <Card key={note.note_uid} className="relative">
                <CardContent className="pt-4">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <h3 className="font-semibold text-lg">{note.title}</h3>
                        {note.version > 1 && (
                          <Badge variant="outline" className="text-xs">
                            v{note.version}
                          </Badge>
                        )}
                      </div>
                      <div className="flex items-center gap-4 text-xs text-gray-500 mb-3">
                        <div className="flex items-center gap-1">
                          <Calendar className="w-3 h-3" />
                          <span>Updated {formatDate(note.updated_at || note.created_at)}</span>
                        </div>
                        {note.active_collaborators && note.active_collaborators > 0 && (
                          <div className="flex items-center gap-1">
                            <Users className="w-3 h-3" />
                            <span>{note.active_collaborators} active</span>
                          </div>
                        )}
                      </div>
                      {collaborators[note.note_uid] && collaborators[note.note_uid].length > 0 && (
                        <CollaboratorAvatars 
                          collaborators={collaborators[note.note_uid]} 
                          size="sm"
                          className="mb-3"
                        />
                      )}
                    </div>
                    
                    {!readOnly && (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreVertical className="w-4 h-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem 
                            onClick={() => handleDeleteNote(note.note_uid)}
                            className="text-red-600"
                          >
                            <Trash2 className="w-4 h-4 mr-2" />
                            Delete Note
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )}
                  </div>
                  
                  <CollaborativeEditor
                    documentId={getDocumentId(note.note_uid)}
                    userName={user?.name || 'Anonymous'}
                    userColor={user ? getUserColor(user.user_uid) : '#3B82F6'}
                    onSave={handleSaveNote(note.note_uid)}
                    autoSave={true}
                    autoSaveInterval={3000}
                    placeholder="Start typing your note..."
                    readOnly={readOnly}
                    className="mt-2"
                  />
                </CardContent>
              </Card>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
};
