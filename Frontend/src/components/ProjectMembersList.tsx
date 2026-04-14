import { useState, useEffect } from 'react';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Users, Crown, Loader2, PenLine } from 'lucide-react';
import { apiService } from '@/services/api';
import { ProjectMember } from '@/types';
import { useToast } from '@/hooks/use-toast';
import { AddProjectMemberDialog } from './AddProjectMemberDialog';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

interface ProjectMembersListProps {
  projectId: string;
}

export function ProjectMembersList({ projectId }: ProjectMembersListProps) {
  const [members, setMembers] = useState<ProjectMember[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { toast } = useToast();

  const loadMembers = async () => {
    try {
      setIsLoading(true);
      const response = await apiService.getProjectMembers(projectId);
      setMembers(response.data || []);
    } catch (error) {
      console.error('Failed to load members:', error);
      toast({
        title: 'Error',
        description: 'Failed to load project members',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadMembers();
  }, [projectId]);

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-muted-foreground">
            {members.length} {members.length === 1 ? 'member' : 'members'}
          </p>
        </div>
        <AddProjectMemberDialog projectId={projectId} onMemberAdded={loadMembers} members={members} />
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : members.length === 0 ? (
        <div className="text-center py-8">
          <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <p className="text-sm text-muted-foreground">No members yet</p>
        </div>
      ) : (
        <TooltipProvider>
          <div className="space-y-2 max-h-[300px] overflow-y-auto pr-2">
            {members.map((member) => (
              <div
                key={member.user_uid}
                className="flex items-center justify-between p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="bg-primary/10 text-primary text-xs">
                      {getInitials(member.name)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="flex items-center gap-2">
                      <p className="font-medium text-sm">{member.name}</p>
                      {member.role === 'owner' && (
                        <Tooltip>
                          <TooltipTrigger>
                            <Crown className="h-3 w-3 text-yellow-500" />
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Project Owner</p>
                          </TooltipContent>
                        </Tooltip>
                      )}
                      {member.can_write && member.role !== 'owner' && (
                        <Tooltip>
                          <TooltipTrigger>
                            <PenLine className="h-3 w-3 text-blue-500" />
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Can Edit</p>
                          </TooltipContent>
                        </Tooltip>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground">{member.email}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={member.role === 'owner' ? 'default' : 'secondary'} className="text-xs">
                    {member.role}
                  </Badge>
                  {member.can_write && member.role !== 'owner' && (
                    <Badge variant="outline" className="text-xs text-blue-600 border-blue-300">
                      Can Edit
                    </Badge>
                  )}
                </div>
              </div>
            ))}
          </div>
        </TooltipProvider>
      )}
    </div>
  );
}
