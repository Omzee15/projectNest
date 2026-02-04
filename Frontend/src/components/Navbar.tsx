import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Search, Plus, Bell, Settings, LogOut, User, Menu, X } from 'lucide-react';
import { CreateProjectDialog } from './CreateProjectDialog';
import { useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';

interface User {
  uid: string;
  email: string;
  name: string;
}

export function Navbar() {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [showSearch, setShowSearch] = useState(false);

  useEffect(() => {
    // Check if user is authenticated
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (token && userData) {
      try {
        const parsedUser = JSON.parse(userData);
        setUser(parsedUser);
        setIsAuthenticated(true);
      } catch (error) {
        console.error('Error parsing user data:', error);
        handleLogout();
      }
    }
  }, []);

  const handleLogout = () => {
    // Clear authentication data
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    setIsAuthenticated(false);
    
    // Redirect to login page
    navigate('/login');
  };

  const getUserInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word.charAt(0))
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <nav className="border-b border-border bg-background h-14 flex items-center justify-between px-4 md:px-6 relative">
      <div className="flex items-center gap-3 md:gap-6 flex-1">
        <button 
          onClick={() => navigate('/')}
          className="font-bold text-base md:text-lg text-foreground hover:text-primary transition-colors whitespace-nowrap"
        >
          ProjectBoard
        </button>
        
        {/* Desktop Navigation */}
        <div className="hidden md:flex items-center gap-4">
          <Button 
            variant="ghost" 
            size="sm"
            onClick={() => navigate('/projects')}
            className="text-muted-foreground hover:text-foreground"
          >
            Projects
          </Button>
        </div>
        
        {/* Desktop Search */}
        <div className="hidden lg:block relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
          <Input
            placeholder="Search..."
            className="pl-10 w-48 xl:w-64 bg-muted/50 border-border-light"
          />
        </div>
        
        {/* Mobile Search Toggle */}
        <Button 
          variant="ghost" 
          size="sm"
          onClick={() => setShowSearch(!showSearch)}
          className="lg:hidden ml-auto"
        >
          <Search className="h-4 w-4" />
        </Button>
      </div>

      <div className="flex items-center gap-2 md:gap-3">
        {isAuthenticated ? (
          <>
            {/* Desktop Create Button */}
            <div className="hidden md:block">
              <CreateProjectDialog 
                trigger={
                  <Button variant="ghost" size="sm">
                    <Plus className="h-4 w-4 mr-2" />
                    Create
                  </Button>
                }
              />
            </div>
            
            {/* Mobile Create Button - Icon Only */}
            <div className="md:hidden">
              <CreateProjectDialog 
                trigger={
                  <Button variant="ghost" size="sm">
                    <Plus className="h-4 w-4" />
                  </Button>
                }
              />
            </div>
            
            <Button variant="ghost" size="sm" className="hidden sm:flex">
              <Bell className="h-4 w-4" />
            </Button>
            
            <Button variant="ghost" size="sm" onClick={() => navigate('/settings')} className="hidden sm:flex">
              <Settings className="h-4 w-4" />
            </Button>
            
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="text-xs bg-primary text-primary-foreground">
                      {user ? getUserInitials(user.name) : 'U'}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">
                      {user?.name || 'User'}
                    </p>
                    <p className="text-xs leading-none text-muted-foreground">
                      {user?.email || ''}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => navigate('/profile')} disabled>
                  <User className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigate('/settings')}>
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            
            {/* Mobile Menu Button */}
            <Button 
              variant="ghost" 
              size="sm" 
              className="md:hidden"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
            </Button>
          </>
        ) : (
          <Button onClick={() => navigate('/login')} size="sm">
            Sign In
          </Button>
        )}
      </div>
      
      {/* Mobile Search Overlay */}
      {showSearch && (
        <div className="absolute top-14 left-0 right-0 bg-background border-b border-border p-4 lg:hidden z-50">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search..."
              className="pl-10 w-full bg-muted/50 border-border-light"
              autoFocus
            />
          </div>
        </div>
      )}
      
      {/* Mobile Menu Dropdown */}
      {mobileMenuOpen && isAuthenticated && (
        <div className="absolute top-14 left-0 right-0 bg-background border-b border-border md:hidden z-50">
          <div className="flex flex-col p-2">
            <Button 
              variant="ghost" 
              className="justify-start"
              onClick={() => {
                navigate('/projects');
                setMobileMenuOpen(false);
              }}
            >
              Projects
            </Button>
            <Button 
              variant="ghost" 
              className="justify-start"
              onClick={() => {
                navigate('/settings');
                setMobileMenuOpen(false);
              }}
            >
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Button>
            <Button 
              variant="ghost" 
              className="justify-start"
            >
              <Bell className="mr-2 h-4 w-4" />
              Notifications
            </Button>
          </div>
        </div>
      )}
    </nav>
  );
}