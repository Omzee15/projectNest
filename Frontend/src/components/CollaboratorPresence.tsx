import React from 'react';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Badge } from '@/components/ui/badge';
import { Users } from 'lucide-react';
import { cn } from '@/lib/utils';

export interface CollaboratorPresence {
  user_uid: string;
  user_name: string;
  user_email?: string;
  color: string;
  cursor_position?: any;
  last_seen_at: string;
}

interface CollaboratorAvatarsProps {
  collaborators: CollaboratorPresence[];
  maxDisplay?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const CollaboratorAvatars: React.FC<CollaboratorAvatarsProps> = ({
  collaborators,
  maxDisplay = 5,
  size = 'md',
  className,
}) => {
  const displayCollaborators = collaborators.slice(0, maxDisplay);
  const remainingCount = Math.max(0, collaborators.length - maxDisplay);

  const sizeClasses = {
    sm: 'w-6 h-6 text-xs',
    md: 'w-8 h-8 text-sm',
    lg: 'w-10 h-10 text-base',
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const getTimeSince = (timestamp: string) => {
    const now = new Date();
    const lastSeen = new Date(timestamp);
    const diffMs = now.getTime() - lastSeen.getTime();
    const diffSecs = Math.floor(diffMs / 1000);
    
    if (diffSecs < 60) return 'just now';
    if (diffSecs < 3600) return `${Math.floor(diffSecs / 60)}m ago`;
    if (diffSecs < 86400) return `${Math.floor(diffSecs / 3600)}h ago`;
    return `${Math.floor(diffSecs / 86400)}d ago`;
  };

  if (collaborators.length === 0) {
    return null;
  }

  return (
    <TooltipProvider>
      <div className={cn('flex items-center gap-1', className)}>
        <Users className="w-4 h-4 text-gray-500 mr-1" />
        <div className="flex -space-x-2">
          {displayCollaborators.map((collaborator) => (
            <Tooltip key={collaborator.user_uid}>
              <TooltipTrigger>
                <Avatar
                  className={cn(
                    sizeClasses[size],
                    'border-2 border-white transition-transform hover:scale-110'
                  )}
                  style={{ 
                    borderColor: collaborator.color,
                    boxShadow: `0 0 0 2px ${collaborator.color}`
                  }}
                >
                  <AvatarFallback
                    style={{ backgroundColor: collaborator.color }}
                    className="text-white font-semibold"
                  >
                    {getInitials(collaborator.user_name)}
                  </AvatarFallback>
                </Avatar>
              </TooltipTrigger>
              <TooltipContent>
                <div className="text-sm">
                  <p className="font-semibold">{collaborator.user_name}</p>
                  {collaborator.user_email && (
                    <p className="text-gray-500 text-xs">{collaborator.user_email}</p>
                  )}
                  <p className="text-gray-400 text-xs mt-1">
                    Active {getTimeSince(collaborator.last_seen_at)}
                  </p>
                </div>
              </TooltipContent>
            </Tooltip>
          ))}
          {remainingCount > 0 && (
            <Tooltip>
              <TooltipTrigger>
                <Avatar className={cn(sizeClasses[size], 'border-2 border-white bg-gray-200')}>
                  <AvatarFallback className="text-gray-600 font-semibold">
                    +{remainingCount}
                  </AvatarFallback>
                </Avatar>
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-sm">{remainingCount} more collaborator{remainingCount > 1 ? 's' : ''}</p>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
        <Badge variant="secondary" className="ml-2">
          {collaborators.length} {collaborators.length === 1 ? 'user' : 'users'} online
        </Badge>
      </div>
    </TooltipProvider>
  );
};

interface ConnectionStatusProps {
  status: 'connected' | 'connecting' | 'disconnected';
  className?: string;
}

export const ConnectionStatus: React.FC<ConnectionStatusProps> = ({ status, className }) => {
  const statusConfig = {
    connected: {
      color: 'bg-green-500',
      text: 'Connected',
      pulse: true,
    },
    connecting: {
      color: 'bg-yellow-500',
      text: 'Connecting...',
      pulse: true,
    },
    disconnected: {
      color: 'bg-red-500',
      text: 'Disconnected',
      pulse: false,
    },
  };

  const config = statusConfig[status];

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <div className="relative">
        <div className={cn('w-2 h-2 rounded-full', config.color)} />
        {config.pulse && (
          <div
            className={cn(
              'absolute inset-0 w-2 h-2 rounded-full animate-ping',
              config.color,
              'opacity-75'
            )}
          />
        )}
      </div>
      <span className="text-xs text-gray-600">{config.text}</span>
    </div>
  );
};
