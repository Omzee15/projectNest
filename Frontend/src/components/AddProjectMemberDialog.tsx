import { useState, useEffect, useRef, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { UserPlus, Loader2, Search, Check } from 'lucide-react';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { UserSearchResult } from '@/types';

interface AddProjectMemberDialogProps {
  projectId: string;
  onMemberAdded?: () => void;
  trigger?: React.ReactNode;
}

export function AddProjectMemberDialog({
  projectId,
  onMemberAdded,
  trigger,
}: AddProjectMemberDialogProps) {
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [users, setUsers] = useState<UserSearchResult[]>([]);
  const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const [isAdding, setIsAdding] = useState(false);
  const { toast } = useToast();
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  const searchUsers = useCallback(async (query: string) => {
    try {
      setIsSearching(true);
      const response = await apiService.searchUsers(query, projectId);
      setUsers(response.data || []);
    } catch (error) {
      console.error('Failed to search users:', error);
      setUsers([]);
    } finally {
      setIsSearching(false);
    }
  }, [projectId]);

  // Load initial user list when dialog opens
  useEffect(() => {
    if (open) {
      searchUsers('');
      setSearchQuery('');
      setSelectedUser(null);
    }
  }, [open, searchUsers]);

  // Debounced search
  useEffect(() => {
    if (!open) return;

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      searchUsers(searchQuery);
    }, 300);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [searchQuery, open, searchUsers]);

  const handleAddMember = async () => {
    if (!selectedUser) {
      toast({
        title: 'No user selected',
        description: 'Please select a user to add',
        variant: 'destructive',
      });
      return;
    }

    try {
      setIsAdding(true);
      await apiService.addProjectMember(projectId, selectedUser.email, 'member');

      toast({
        title: 'Member added!',
        description: `${selectedUser.name} has been added to the project`,
      });

      setSearchQuery('');
      setSelectedUser(null);
      setOpen(false);

      if (onMemberAdded) {
        onMemberAdded();
      }
    } catch (error: any) {
      console.error('Failed to add member:', error);

      let errorMessage = 'Failed to add member. Please try again.';
      if (error.message?.includes('not found')) {
        errorMessage = 'User not found. They need to register first.';
      } else if (error.message?.includes('already a member')) {
        errorMessage = 'This user is already a member of the project.';
      }

      toast({
        title: 'Error',
        description: errorMessage,
        variant: 'destructive',
      });
    } finally {
      setIsAdding(false);
    }
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <UserPlus className="mr-2 h-4 w-4" />
            Add Member
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Project Member</DialogTitle>
          <DialogDescription>
            Search and select a user to add to this project.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-3 py-2">
          {/* Search input */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by name or email..."
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value);
                setSelectedUser(null);
              }}
              className="pl-9"
              autoFocus
            />
          </div>

          {/* User list */}
          <div className="border rounded-lg max-h-[280px] min-h-[120px] overflow-y-auto">
            {isSearching ? (
              <div className="flex items-center justify-center py-8">
                <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                <span className="ml-2 text-sm text-muted-foreground">Searching...</span>
              </div>
            ) : users.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
                <Search className="h-8 w-8 mb-2 opacity-50" />
                <p className="text-sm">
                  {searchQuery ? 'No users found' : 'No users available to add'}
                </p>
                {searchQuery && (
                  <p className="text-xs mt-1">Try a different search term</p>
                )}
              </div>
            ) : (
              <div className="divide-y">
                {users.map((user) => {
                  const isSelected = selectedUser?.user_uid === user.user_uid;
                  return (
                    <button
                      key={user.user_uid}
                      type="button"
                      onClick={() => setSelectedUser(isSelected ? null : user)}
                      className={`w-full flex items-center gap-3 px-3 py-2.5 text-left transition-colors hover:bg-accent/50 ${
                        isSelected ? 'bg-primary/10 hover:bg-primary/15' : ''
                      }`}
                    >
                      <Avatar className="h-8 w-8 shrink-0">
                        <AvatarFallback className="bg-primary/10 text-primary text-xs">
                          {getInitials(user.name)}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-sm truncate">{user.name}</p>
                        <p className="text-xs text-muted-foreground truncate">{user.email}</p>
                      </div>
                      {isSelected && (
                        <Check className="h-4 w-4 text-primary shrink-0" />
                      )}
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </div>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => setOpen(false)}
            disabled={isAdding}
          >
            Cancel
          </Button>
          <Button
            type="button"
            onClick={handleAddMember}
            disabled={isAdding || !selectedUser}
          >
            {isAdding && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Add {selectedUser ? selectedUser.name : 'Member'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
