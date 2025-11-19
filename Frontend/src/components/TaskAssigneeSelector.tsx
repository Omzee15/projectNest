import { useState, useEffect } from 'react';
import { Check, ChevronsUpDown, Users, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from '@/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { ProjectMember } from '@/types';
import { apiService } from '@/services/api';
import { useToast } from '@/hooks/use-toast';
import { cn } from '@/lib/utils';

interface TaskAssigneeSelectorProps {
  projectId: string;
  selectedUserIds: string[];
  onSelectionChange: (userIds: string[]) => void;
  disabled?: boolean;
}

export function TaskAssigneeSelector({
  projectId,
  selectedUserIds,
  onSelectionChange,
  disabled = false,
}: TaskAssigneeSelectorProps) {
  const [open, setOpen] = useState(false);
  const [members, setMembers] = useState<ProjectMember[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    loadMembers();
  }, [projectId]);

  const loadMembers = async () => {
    try {
      setIsLoading(true);
      const response = await apiService.getProjectMembers(projectId);
      setMembers(response.data || []);
    } catch (error) {
      console.error('Failed to load project members:', error);
      toast({
        title: 'Error',
        description: 'Failed to load project members',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const toggleMember = (userUid: string) => {
    const newSelection = selectedUserIds.includes(userUid)
      ? selectedUserIds.filter((id) => id !== userUid)
      : [...selectedUserIds, userUid];
    onSelectionChange(newSelection);
  };

  const removeMember = (userUid: string, e: React.MouseEvent) => {
    e.stopPropagation();
    onSelectionChange(selectedUserIds.filter((id) => id !== userUid));
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const selectedMembers = members.filter((m) => selectedUserIds.includes(m.user_uid));

  return (
    <div className="space-y-2">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
            disabled={disabled}
          >
            <div className="flex items-center gap-2">
              <Users className="h-4 w-4 opacity-50" />
              <span>
                {selectedUserIds.length === 0
                  ? 'Assign to...'
                  : `${selectedUserIds.length} member${selectedUserIds.length > 1 ? 's' : ''} assigned`}
              </span>
            </div>
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-full p-0" align="start">
          <Command>
            <CommandInput placeholder="Search members..." />
            <CommandEmpty>
              {isLoading ? 'Loading members...' : 'No members found.'}
            </CommandEmpty>
            <CommandGroup className="max-h-64 overflow-auto">
              {members.map((member) => (
                <CommandItem
                  key={member.user_uid}
                  value={member.name}
                  onSelect={() => toggleMember(member.user_uid)}
                  className="cursor-pointer"
                >
                  <Check
                    className={cn(
                      'mr-2 h-4 w-4',
                      selectedUserIds.includes(member.user_uid) ? 'opacity-100' : 'opacity-0'
                    )}
                  />
                  <Avatar className="h-6 w-6 mr-2">
                    <AvatarFallback className="text-xs bg-primary/10 text-primary">
                      {getInitials(member.name)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col">
                    <span className="text-sm">{member.name}</span>
                    <span className="text-xs text-muted-foreground">{member.email}</span>
                  </div>
                </CommandItem>
              ))}
            </CommandGroup>
          </Command>
        </PopoverContent>
      </Popover>

      {selectedMembers.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedMembers.map((member) => (
            <Badge key={member.user_uid} variant="secondary" className="pl-2 pr-1">
              <Avatar className="h-4 w-4 mr-1">
                <AvatarFallback className="text-[8px] bg-primary/10 text-primary">
                  {getInitials(member.name)}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs">{member.name}</span>
              <button
                type="button"
                onClick={(e) => removeMember(member.user_uid, e)}
                className="ml-1 hover:bg-secondary-foreground/20 rounded-full p-0.5"
                disabled={disabled}
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}
