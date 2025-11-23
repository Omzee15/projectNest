import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { ProjectNotesCollaborative } from '@/components/ProjectNotesCollaborative';
import { apiService } from '@/services/api';
import { toast } from '@/hooks/use-toast';
import { Note, NoteRequest, NoteUpdateRequest } from '@/types';

const ProjectNotesNotion: React.FC = () => {
  const { projectUid } = useParams<{ projectUid: string }>();
  const queryClient = useQueryClient();
  const [draggedNote, setDraggedNote] = useState<Note | null>(null);
  const [editingFolder, setEditingFolder] = useState<NoteFolder | null>(null);
  const [showEditFolderDialog, setShowEditFolderDialog] = useState(false);
  const [editFolderName, setEditFolderName] = useState('');

  // Fetch project
  const { data: project } = useQuery({
    queryKey: ['project', projectUid],
    queryFn: () => apiService.getProject(projectUid!),
    enabled: !!projectUid,
  });

  // Fetch notes
  const { data: notesResponse, isLoading: isNotesLoading } = useQuery({
    queryKey: ['notes', projectUid],
    queryFn: () => apiService.getNotesByProject(projectUid!),
    enabled: !!projectUid,
  });

  // Fetch folders
  const { data: foldersResponse, isLoading: isFoldersLoading } = useQuery({
    queryKey: ['folders', projectUid],
    queryFn: () => apiService.getFolders(projectUid!),
    enabled: !!projectUid,
  });

  // Create note mutation
  const createNoteMutation = useMutation({
    mutationFn: (noteData: NoteRequest) => 
      apiService.createNote(projectUid!, noteData),
    onSuccess: () => {
      toast({
        title: 'Note created',
        description: 'Your note has been created successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      setNewNoteName('');
      setShowCreateNoteDialog(false);
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to create note. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Create folder mutation
  const createFolderMutation = useMutation({
    mutationFn: (folderData: NoteFolderRequest) => 
      apiService.createFolder(projectUid!, folderData),
    onSuccess: () => {
      toast({
        title: 'Folder created',
        description: 'Your folder has been created successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['folders', projectUid] });
      setNewFolderName('');
      setShowCreateFolderDialog(false);
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to create folder. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Update note mutation
  const updateNoteMutation = useMutation({
    mutationFn: ({ noteUid, noteData }: { noteUid: string; noteData: NoteRequest }) => 
      apiService.updateNote(noteUid, noteData),
    onSuccess: () => {
      toast({
        title: 'Note saved',
        description: 'Your changes have been saved.',
      });
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to save note. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Delete note mutation
  const deleteNoteMutation = useMutation({
    mutationFn: (noteUid: string) => apiService.deleteNote(noteUid),
    onSuccess: () => {
      toast({
        title: 'Note deleted',
        description: 'Your note has been deleted successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      setSelectedNote(null);
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to delete note. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const moveNoteMutation = useMutation({
    mutationFn: ({ noteUid, folderId }: { noteUid: string; folderId?: number }) => 
      apiService.moveNoteToFolder(noteUid, folderId),
    onSuccess: () => {
      toast({
        title: 'Note moved',
        description: 'Your note has been moved successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      setDraggedNote(null);
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to move note. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Update folder mutation
  const updateFolderMutation = useMutation({
    mutationFn: ({ folderUid, folderData }: { folderUid: string; folderData: NoteFolderRequest }) => 
      apiService.updateFolder(folderUid, folderData),
    onSuccess: () => {
      toast({
        title: 'Folder updated',
        description: 'Your folder has been updated successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['folders', projectUid] });
      setShowEditFolderDialog(false);
      setEditingFolder(null);
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to update folder. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Delete folder mutation
  const deleteFolderMutation = useMutation({
    mutationFn: (folderUid: string) => apiService.deleteFolder(folderUid),
    onSuccess: () => {
      toast({
        title: 'Folder deleted',
        description: 'Your folder has been deleted successfully.',
      });
      queryClient.invalidateQueries({ queryKey: ['folders', projectUid] });
      queryClient.invalidateQueries({ queryKey: ['notes', projectUid] }); // Refresh notes as they might have been moved
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to delete folder. Please try again.',
        variant: 'destructive',
      });
    },
  });

  // Load selected note content
  useEffect(() => {
    if (selectedNote) {
      setEditingTitle(selectedNote.title);
      // Extract text content from blocks
      const textContent = selectedNote.content.blocks
        .filter(block => block.type === 'text' && block.content)
        .map(block => block.content)
        .join('\\n\\n');
      setEditingContent(textContent);
      
      // Mark auto-save as clean for the newly loaded note
      setTimeout(() => {
        autoSaveContent.markAsSaved();
        autoSaveTitle.markAsSaved();
      }, 100);
    }
  }, [selectedNote]);

  // Auto-save functionality for note content
  const autoSaveContent = useAutoSave(
    editingContent,
    async (content: string) => {
      if (!selectedNote || !editingTitle.trim()) return;

      const noteContent: NoteContent = {
        blocks: [{
          id: Date.now().toString(),
          type: 'text',
          content: content
        }]
      };

      await apiService.updateNote(selectedNote.note_uid, {
        title: editingTitle,
        content: noteContent,
      });
    },
    {
      delay: 2000, // 2 seconds after user stops typing
      enabled: !!selectedNote && !!editingTitle.trim(),
      onSaveStart: () => {
        // Optional: Show saving indicator
      },
      onSaveSuccess: () => {
        toast({
          title: 'Note saved',
          description: 'Your changes have been automatically saved.',
          duration: 2000,
        });
        // Invalidate queries to refresh the note list
        queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      },
      onSaveError: (error) => {
        toast({
          title: 'Error',
          description: 'Failed to auto-save note. Please try saving manually.',
          variant: 'destructive',
        });
      },
    }
  );

  // Auto-save functionality for note title
  const autoSaveTitle = useAutoSave(
    editingTitle,
    async (title: string) => {
      if (!selectedNote || !title.trim()) return;

      const noteContent: NoteContent = {
        blocks: [{
          id: Date.now().toString(),
          type: 'text',
          content: editingContent
        }]
      };

      await apiService.updateNote(selectedNote.note_uid, {
        title: title,
        content: noteContent,
      });
    },
    {
      delay: 1500, // 1.5 seconds for title changes
      enabled: !!selectedNote,
      onSaveSuccess: () => {
        // Refresh the note list to show updated title
        queryClient.invalidateQueries({ queryKey: ['notes', projectUid] });
      },
      onSaveError: (error) => {
        toast({
          title: 'Error',
          description: 'Failed to auto-save note title.',
          variant: 'destructive',
        });
      },
    }
  );

  const handleCreateNote = () => {
    if (!newNoteName.trim()) return;

    const content: NoteContent = {
      blocks: [{
        id: Date.now().toString(),
        type: 'text',
        content: ''
      }]
    };

    // Find the folder's internal ID based on the UID
    let folderId: number | undefined = undefined;
    if (createNoteInFolder) {
      const folder = folders.find(f => f.folder_uid === createNoteInFolder);
      folderId = folder ? folder.id : undefined; // Use the folder's internal database ID
    }

    createNoteMutation.mutate({
      title: newNoteName,
      content: content,
      folder_id: folderId,
    });
  };

  const handleCreateFolder = () => {
    if (!newFolderName.trim()) return;

    createFolderMutation.mutate({
      name: newFolderName,
      parent_folder_id: selectedFolder ? parseInt(selectedFolder) : undefined,
    });
  };

  const handleSaveNote = async () => {
    if (!selectedNote || !editingTitle.trim()) return;

    try {
      // Use the auto-save's manual save function
      await autoSaveContent.save();
      await autoSaveTitle.save();
      
      toast({
        title: 'Note saved',
        description: 'Your note has been saved successfully.',
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to save note. Please try again.',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteNote = () => {
    if (!selectedNote) return;
    if (window.confirm('Are you sure you want to delete this note?')) {
      deleteNoteMutation.mutate(selectedNote.note_uid);
    }
  };

  const toggleFolder = (folderId: string) => {
    const newExpanded = new Set(expandedFolders);
    if (newExpanded.has(folderId)) {
      newExpanded.delete(folderId);
    } else {
      newExpanded.add(folderId);
    }
    setExpandedFolders(newExpanded);
  };

  // Drag and drop handlers
  const handleDragStart = (note: Note) => {
    setDraggedNote(note);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDropOnFolder = (e: React.DragEvent, folderUid: string) => {
    e.preventDefault();
    if (draggedNote) {
      const folder = folders.find(f => f.folder_uid === folderUid);
      if (folder) {
        moveNoteMutation.mutate({
          noteUid: draggedNote.note_uid,
          folderId: folder.id, // Use the folder's internal database ID
        });
      }
    }
  };

  const handleDropOnRoot = (e: React.DragEvent) => {
    e.preventDefault();
    if (draggedNote) {
      moveNoteMutation.mutate({
        noteUid: draggedNote.note_uid,
        folderId: undefined, // Move to root (no folder)
      });
    }
  };

  const handleCreateNoteInFolder = (folderId: string) => {
    setCreateNoteInFolder(folderId);
    setNewNoteName('');
    setShowCreateNoteDialog(true);
  };

  const handleEditFolder = (folder: NoteFolder) => {
    setEditingFolder(folder);
    setEditFolderName(folder.name);
    setShowEditFolderDialog(true);
  };

  const handleUpdateFolder = () => {
    if (!editingFolder || !editFolderName.trim()) return;

    updateFolderMutation.mutate({
      folderUid: editingFolder.folder_uid,
      folderData: {
        name: editFolderName.trim(),
        parent_folder_id: editingFolder.parent_folder_id,
        position: editingFolder.position,
      },
    });
  };

  const handleDeleteFolder = (folder: NoteFolder) => {
    if (window.confirm(`Are you sure you want to delete the folder "${folder.name}"? All notes in this folder will be moved to the root.`)) {
      deleteFolderMutation.mutate(folder.folder_uid);
    }
  };

  if (!projectUid) {
    return <div>Project not found</div>;
  }

  const notes = notesResponse?.data?.notes || [];
  const folders = foldersResponse?.data?.folders || [];

  return (
    <div className="h-screen flex flex-col bg-gray-50 overflow-hidden">
      {/* Header */}
      <div className="bg-white border-b px-4 py-3 flex items-center gap-4 flex-shrink-0">
        <Button 
          variant="outline" 
          size="sm"
          onClick={() => window.history.back()}
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Project
        </Button>
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded bg-gradient-to-br from-green-500 to-blue-500">
            <FileText className="w-4 h-4 text-white" />
          </div>
          <div>
            <h1 className="font-semibold">{project?.data?.name || 'Project'} Notes</h1>
            <p className="text-sm text-gray-600">Organize your thoughts and ideas</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="notes" className="flex flex-1 flex-col overflow-hidden">
        <div className="bg-white flex-shrink-0 h-10 flex items-center px-4">
          <TabsList className="border-0 bg-transparent p-0">
            <TabsTrigger value="notes" className="flex items-center gap-2 h-9">
              <FileText className="w-4 h-4" />
              Notes
            </TabsTrigger>
            <TabsTrigger value="diff" className="flex items-center gap-2 h-9">
              <GitCompare className="w-4 h-4" />
              Diff Checker
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="notes" className="flex flex-1 overflow-hidden !mt-0 m-0 p-0">
          <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <div className="w-80 bg-white border-r flex flex-col">
          {/* Sidebar Header */}
          <div className="p-4 border-b">
            <div className="flex items-center justify-between mb-3">
              <h2 className="font-medium">Notes & Folders</h2>
              <div className="flex gap-1">
                <Dialog open={showCreateFolderDialog} onOpenChange={setShowCreateFolderDialog}>
                  <DialogTrigger asChild>
                    <Button variant="ghost" size="sm">
                      <FolderPlus className="w-4 h-4" />
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Create New Folder</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Input
                        value={newFolderName}
                        onChange={(e) => setNewFolderName(e.target.value)}
                        placeholder="Folder name"
                        onKeyDown={(e) => e.key === 'Enter' && handleCreateFolder()}
                      />
                      <Button onClick={handleCreateFolder} disabled={!newFolderName.trim()}>
                        Create Folder
                      </Button>
                    </div>
                  </DialogContent>
                </Dialog>

                {/* Edit Folder Dialog */}
                <Dialog open={showEditFolderDialog} onOpenChange={setShowEditFolderDialog}>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Edit Folder</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Input
                        value={editFolderName}
                        onChange={(e) => setEditFolderName(e.target.value)}
                        placeholder="Folder name"
                        onKeyDown={(e) => e.key === 'Enter' && handleUpdateFolder()}
                      />
                      <div className="flex gap-2">
                        <Button onClick={handleUpdateFolder} disabled={!editFolderName.trim()}>
                          Update Folder
                        </Button>
                        <Button variant="outline" onClick={() => setShowEditFolderDialog(false)}>
                          Cancel
                        </Button>
                      </div>
                    </div>
                  </DialogContent>
                </Dialog>
                
                <Dialog open={showCreateNoteDialog} onOpenChange={(open) => {
                  setShowCreateNoteDialog(open);
                  if (!open) {
                    setCreateNoteInFolder(null);
                    setNewNoteName('');
                  }
                }}>
                  <DialogTrigger asChild>
                    <Button variant="ghost" size="sm">
                      <Plus className="w-4 h-4" />
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>
                        Create New Note
                        {createNoteInFolder && (
                          <span className="text-sm font-normal text-gray-600">
                            {' '}in {folders.find(f => f.folder_uid === createNoteInFolder)?.name}
                          </span>
                        )}
                      </DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Input
                        value={newNoteName}
                        onChange={(e) => setNewNoteName(e.target.value)}
                        placeholder="Note title"
                        onKeyDown={(e) => e.key === 'Enter' && handleCreateNote()}
                      />
                      <div className="flex gap-2">
                        <Button onClick={handleCreateNote} disabled={!newNoteName.trim()}>
                          Create Note
                        </Button>
                        {createNoteInFolder && (
                          <Button 
                            variant="outline" 
                            onClick={() => {
                              setCreateNoteInFolder(null);
                            }}
                          >
                            Create in Root
                          </Button>
                        )}
                      </div>
                    </div>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </div>

          {/* Sidebar Content */}
          <div className="flex-1 overflow-y-auto p-2">
            {isNotesLoading || isFoldersLoading ? (
              <div className="space-y-2">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-8 bg-gray-200 rounded animate-pulse" />
                ))}
              </div>
            ) : (
              <div 
                className="space-y-2"
                onDragOver={handleDragOver}
                onDrop={handleDropOnRoot}
              >
                {/* Folders Section */}
                {folders.filter(f => !f.parent_folder_id).length > 0 && (
                  <div className="text-xs font-medium text-gray-500 uppercase tracking-wide px-3 py-1">
                    Folders
                  </div>
                )}
                
                {/* Root level folders */}
                {folders.filter(f => !f.parent_folder_id).map((folder) => (
                  <div key={folder.folder_uid}>
                    <div 
                      className="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-blue-50 cursor-pointer text-sm group font-medium border border-transparent hover:border-blue-200 transition-all"
                      onClick={() => toggleFolder(folder.folder_uid)}
                      onDragOver={handleDragOver}
                      onDrop={(e) => handleDropOnFolder(e, folder.folder_uid)}
                    >
                      {expandedFolders.has(folder.folder_uid) ? 
                        <ChevronDown className="w-4 h-4 text-gray-600" /> : 
                        <ChevronRight className="w-4 h-4 text-gray-600" />
                      }
                      <Folder className="w-5 h-5 text-blue-600" />
                      <span className="flex-1 text-gray-800">{folder.name}</span>
                      <div className="text-xs text-gray-500 mr-2">
                        {notes.filter(note => note.folder_id === folder.id).length} notes
                      </div>
                      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100">
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0 hover:bg-blue-100"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleCreateNoteInFolder(folder.folder_uid);
                          }}
                          title="Add note to folder"
                        >
                          <Plus className="w-4 h-4 text-blue-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0 hover:bg-yellow-100"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEditFolder(folder);
                          }}
                          title="Edit folder"
                        >
                          <Edit2 className="w-4 h-4 text-yellow-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0 hover:bg-red-100"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteFolder(folder);
                          }}
                          title="Delete folder"
                        >
                          <Trash2 className="w-4 h-4 text-red-600" />
                        </Button>
                      </div>
                    </div>
                    
                    {/* Notes within folder */}
                    {expandedFolders.has(folder.folder_uid) && (
                      <div className="ml-6 border-l-2 border-gray-100 pl-2 space-y-1">
                        {notes.filter(note => note.folder_id === folder.id).map((note) => (
                          <div 
                            key={note.note_uid}
                            className={`flex items-center gap-2 px-3 py-2 rounded-md cursor-pointer text-sm transition-colors ${
                              selectedNote?.note_uid === note.note_uid 
                                ? 'bg-blue-50 border-l-2 border-blue-400' 
                                : 'hover:bg-gray-50 border-l-2 border-transparent'
                            }`}
                            draggable
                            onDragStart={() => handleDragStart(note)}
                            onClick={() => setSelectedNote(note)}
                          >
                            <div className="w-2 h-2 rounded-full bg-gray-300" /> {/* Bullet point */}
                            <FileText className="w-4 h-4 text-gray-400" />
                            <span className="flex-1 truncate text-gray-700">{note.title}</span>
                            <div className="text-xs text-gray-400">
                              {notes.filter(n => n.folder_id === folder.id).length > 1 ? 
                                `${notes.filter(n => n.folder_id === folder.id).indexOf(note) + 1}/${notes.filter(n => n.folder_id === folder.id).length}` 
                                : ''
                              }
                            </div>
                          </div>
                        ))}
                        {notes.filter(note => note.folder_id === folder.id).length === 0 && (
                          <div className="flex items-center gap-2 px-3 py-2 text-sm text-gray-400 italic">
                            <div className="w-2 h-2" /> {/* Empty space for alignment */}
                            <FileText className="w-4 h-4 text-gray-300" />
                            <span>No notes in this folder</span>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                ))}

                {/* Root level notes */}
                <div className="space-y-1">
                  {notes.filter(note => !note.folder_id).length > 0 && (
                    <div className="text-xs font-medium text-gray-500 uppercase tracking-wide px-3 py-1 mt-4">
                      Root Notes
                    </div>
                  )}
                  {notes.filter(note => !note.folder_id).map((note) => (
                    <div 
                      key={note.note_uid}
                      className={`flex items-center gap-2 px-3 py-2 rounded-md cursor-pointer text-sm transition-colors ${
                        selectedNote?.note_uid === note.note_uid 
                          ? 'bg-blue-50 border-l-2 border-blue-400' 
                          : 'hover:bg-gray-50 border-l-2 border-transparent'
                      }`}
                      draggable
                      onDragStart={() => handleDragStart(note)}
                      onClick={() => setSelectedNote(note)}
                    >
                      <FileText className="w-4 h-4 text-gray-500" />
                      <span className="flex-1 truncate text-gray-700">{note.title}</span>
                    </div>
                  ))}
                </div>

                {/* Empty state */}
                {notes.length === 0 && folders.length === 0 && (
                  <div className="text-center py-8 text-gray-500">
                    <FileText className="w-8 h-8 mx-auto mb-2 opacity-50" />
                    <p className="text-sm">No notes or folders yet</p>
                    <p className="text-xs">Create your first note or folder</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col">
          {selectedNote ? (
            <>
              {/* Note Header */}
              <div className="bg-white border-b px-6 py-4 flex items-center justify-between">
                <div className="flex-1 mr-4">
                  <Input
                    value={editingTitle}
                    onChange={(e) => setEditingTitle(e.target.value)}
                    onKeyDown={(e) => {
                      // Handle Ctrl+S / Cmd+S for manual save
                      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
                        e.preventDefault();
                        handleSaveNote();
                      }
                    }}
                    className="text-lg font-semibold border-none px-0 focus:ring-0"
                    placeholder="Untitled"
                  />
                </div>
                <div className="flex items-center gap-3">
                  {/* Auto-save status indicator */}
                  <div className="flex items-center gap-2 text-sm text-gray-500">
                    {(autoSaveContent.isSaving || autoSaveTitle.isSaving) ? (
                      <>
                        <Clock className="w-3 h-3 animate-pulse" />
                        <span>Saving...</span>
                      </>
                    ) : (autoSaveContent.hasUnsavedChanges || autoSaveTitle.hasUnsavedChanges) ? (
                      <>
                        <div className="w-2 h-2 bg-orange-500 rounded-full" />
                        <span>Unsaved changes</span>
                      </>
                    ) : autoSaveContent.lastSaved ? (
                      <>
                        <div className="w-2 h-2 bg-green-500 rounded-full" />
                        <span>
                          Saved {autoSaveContent.lastSaved.toLocaleTimeString([], { 
                            hour: '2-digit', 
                            minute: '2-digit' 
                          })}
                        </span>
                      </>
                    ) : null}
                  </div>
                  
                  <div className="flex gap-2">
                    <Button 
                      variant="outline" 
                      size="sm" 
                      onClick={handleSaveNote}
                      disabled={autoSaveContent.isSaving || autoSaveTitle.isSaving}
                    >
                      <Save className="w-3 h-3 mr-1" />
                      Save Now
                    </Button>
                    <Button variant="ghost" size="sm" onClick={handleDeleteNote}>
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>

              {/* Note Content */}
              <div className="flex-1 p-6 relative">
                <Textarea
                  value={editingContent}
                  onChange={(e) => setEditingContent(e.target.value)}
                  onKeyDown={(e) => {
                    // Handle Ctrl+S / Cmd+S for manual save
                    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
                      e.preventDefault();
                      handleSaveNote();
                    }
                  }}
                  placeholder="Start writing... (Auto-saves as you type, Ctrl+S to save now)"
                  className={`w-full h-full resize-none border-none focus:ring-0 text-base leading-relaxed transition-all ${
                    autoSaveContent.isSaving ? 'opacity-75' : ''
                  }`}
                />
                
                {/* Floating auto-save indicator */}
                {autoSaveContent.isSaving && (
                  <div className="absolute top-4 right-4 bg-blue-50 text-blue-600 px-2 py-1 rounded-md text-xs flex items-center gap-1 shadow-sm">
                    <Clock className="w-3 h-3 animate-pulse" />
                    Auto-saving...
                  </div>
                )}
              </div>
            </>
          ) : (
            /* Empty state */
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">Select a note to edit</h3>
                <p className="text-gray-600 mb-4">Choose a note from the sidebar or create a new one</p>
                <Button onClick={() => setShowCreateNoteDialog(true)}>
                  <Plus className="w-4 h-4 mr-2" />
                  Create New Note
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
        </TabsContent>

        <TabsContent value="diff" className="flex-1 overflow-hidden m-0 p-0 border-t">
          <DiffChecker />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default ProjectNotesNotion;