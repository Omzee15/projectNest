import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select';
import { Switch } from '../components/ui/switch';
import { Label } from '../components/ui/label';
import { Input } from '../components/ui/input';
import { Separator } from '../components/ui/separator';
import { useToast } from '../hooks/use-toast';
import { useUserSettings, UserSettings } from '../hooks/use-user-settings';
import { useTheme } from '../contexts/ThemeContext';
import { ThemePreview } from '../components/ThemePreview';
import { Palette, Globe, Bell, Volume2, Zap, Clock, User, Grid, ArrowLeft } from 'lucide-react';

const Settings: React.FC = () => {
  const navigate = useNavigate();
  const { settings, loading, error, updateSettings, resetToDefaults } = useUserSettings();
  const { availableThemes, setTheme: applyTheme, updateCurrentTheme } = useTheme();
  const [localSettings, setLocalSettings] = useState<UserSettings | null>(null);
  const [saving, setSaving] = useState(false);
  const [resetting, setResetting] = useState(false);
  const { toast } = useToast();

  // Initialize local settings when settings are loaded
  useEffect(() => {
    if (settings) {
      setLocalSettings(settings);
      // Update theme context to match user's saved theme
      if (settings.theme) {
        updateCurrentTheme(settings.theme);
      }
    }
  }, [settings, updateCurrentTheme]);

  const languages = [
    { code: 'en', name: 'English' },
    { code: 'es', name: 'Español' },
    { code: 'fr', name: 'Français' },
    { code: 'de', name: 'Deutsch' },
    { code: 'it', name: 'Italiano' },
    { code: 'pt', name: 'Português' },
    { code: 'ru', name: 'Русский' },
    { code: 'ja', name: '日本語' },
    { code: 'ko', name: '한국어' },
    { code: 'zh', name: '中文' },
  ];

  const timezones = [
    { value: 'UTC', label: 'UTC' },
    { value: 'America/New_York', label: 'Eastern Time' },
    { value: 'America/Chicago', label: 'Central Time' },
    { value: 'America/Denver', label: 'Mountain Time' },
    { value: 'America/Los_Angeles', label: 'Pacific Time' },
    { value: 'Europe/London', label: 'London' },
    { value: 'Europe/Paris', label: 'Paris' },
    { value: 'Europe/Berlin', label: 'Berlin' },
    { value: 'Asia/Tokyo', label: 'Tokyo' },
    { value: 'Asia/Seoul', label: 'Seoul' },
    { value: 'Asia/Shanghai', label: 'Shanghai' },
    { value: 'Australia/Sydney', label: 'Sydney' },
  ];

  const handleSettingsChange = (key: keyof UserSettings, value: any) => {
    if (localSettings) {
      setLocalSettings(prev => prev ? { ...prev, [key]: value } : null);
    }
  };

  const handleThemeChange = async (themeKey: string) => {
    try {
      // Apply theme immediately
      await applyTheme(themeKey);
      
      // Update local state
      handleSettingsChange('theme', themeKey);
      
      // Save to backend immediately
      if (localSettings) {
        const updatedSettings = { ...localSettings, theme: themeKey };
        const success = await updateSettings(updatedSettings);
        
        if (success) {
          toast({
            title: 'Theme Changed',
            description: `Successfully switched to ${themeKey.replace(/-/g, ' ')}`,
          });
        } else {
          throw new Error('Failed to save theme preference');
        }
      }
    } catch (error) {
      console.error('Failed to apply theme:', error);
      toast({
        title: 'Error',
        description: 'Failed to apply theme',
        variant: 'destructive',
      });
    }
  };

  const saveSettings = async () => {
    if (!localSettings) return;
    
    setSaving(true);
    try {
      const success = await updateSettings(localSettings);
      if (success) {
        toast({
          title: 'Success',
          description: 'Settings saved successfully',
        });
      } else {
        throw new Error('Failed to save settings');
      }
    } catch (error) {
      console.error('Error saving settings:', error);
      toast({
        title: 'Error',
        description: 'Failed to save settings',
        variant: 'destructive',
      });
    } finally {
      setSaving(false);
    }
  };

  const handleResetToDefaults = async () => {
    setResetting(true);
    try {
      const success = await resetToDefaults();
      if (success) {
        toast({
          title: 'Success',
          description: 'Settings reset to defaults successfully',
        });
      } else {
        throw new Error('Failed to reset settings');
      }
    } catch (error) {
      console.error('Error resetting settings:', error);
      toast({
        title: 'Error',
        description: 'Failed to reset settings',
        variant: 'destructive',
      });
    } finally {
      setResetting(false);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto p-6">
        <div className="animate-pulse space-y-6">
          <div className="h-8 w-48 bg-muted rounded"></div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="space-y-6">
              <Card className="p-6">
                <div className="h-6 w-32 bg-muted rounded mb-4"></div>
                <div className="space-y-3">
                  <div className="h-4 w-full bg-muted rounded"></div>
                  <div className="h-10 w-full bg-muted rounded"></div>
                </div>
              </Card>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-6">
        <div className="mb-8">
          <div className="flex items-center gap-4 mb-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate(-1)}
              className="flex items-center gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              Back
            </Button>
          </div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <User className="h-8 w-8" />
            Settings
          </h1>
        </div>
        <Card className="p-6">
          <div className="text-center space-y-4">
            <div className="text-destructive">
              <h2 className="text-xl font-semibold">Error Loading Settings</h2>
              <p className="text-muted-foreground mt-2">{error}</p>
            </div>
            <Button onClick={() => window.location.reload()}>
              Retry
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6">
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate(-1)}
            className="flex items-center gap-2"
          >
            <ArrowLeft className="h-4 w-4" />
            Back
          </Button>
        </div>
        <h1 className="text-3xl font-bold flex items-center gap-2">
          <User className="h-8 w-8" />
          Settings
        </h1>
        <p className="text-muted-foreground mt-2">
          Customize your ProjectNest experience
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Appearance Settings */}
        <Card className="p-6">
          <div className="mb-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Palette className="h-5 w-5" />
              Appearance
            </h2>
            <p className="text-sm text-muted-foreground mt-1">
              Customize the look and feel of your workspace
            </p>
          </div>

          <div className="space-y-4">
            <div>
              <Label htmlFor="theme" className="flex items-center gap-2 mb-3">
                <Grid className="h-4 w-4" />
                Theme Selection
              </Label>
              <div className="grid grid-cols-2 gap-3">
                {availableThemes.map((theme) => (
                  <ThemePreview
                    key={theme.key}
                    theme={theme}
                    isSelected={localSettings?.theme === theme.key}
                    onSelect={handleThemeChange}
                  />
                ))}
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="compact-mode">Compact Mode</Label>
                <p className="text-sm text-muted-foreground">
                  Reduce spacing and padding in the interface
                </p>
              </div>
              <Switch
                id="compact-mode"
                checked={localSettings?.compact_mode || false}
                onCheckedChange={(checked) => handleSettingsChange('compact_mode', checked)}
              />
            </div>
          </div>
        </Card>

        {/* Language & Region */}
        <Card className="p-6">
          <div className="mb-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Globe className="h-5 w-5" />
              Language & Region
            </h2>
            <p className="text-sm text-muted-foreground mt-1">
              Set your language and timezone preferences
            </p>
          </div>

          <div className="space-y-4">
            <div>
              <Label htmlFor="language">Language</Label>
              <Select
                value={localSettings?.language || 'en'}
                onValueChange={(value) => handleSettingsChange('language', value)}
              >
                <SelectTrigger className="mt-1">
                  <SelectValue placeholder="Select language" />
                </SelectTrigger>
                <SelectContent>
                  {languages.map((lang) => (
                    <SelectItem key={lang.code} value={lang.code}>
                      {lang.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label htmlFor="timezone">Timezone</Label>
              <Select
                value={localSettings?.timezone || 'UTC'}
                onValueChange={(value) => handleSettingsChange('timezone', value)}
              >
                <SelectTrigger className="mt-1">
                  <SelectValue placeholder="Select timezone" />
                </SelectTrigger>
                <SelectContent>
                  {timezones.map((tz) => (
                    <SelectItem key={tz.value} value={tz.value}>
                      {tz.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </Card>

        {/* Notifications */}
        <Card className="p-6">
          <div className="mb-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Bell className="h-5 w-5" />
              Notifications
            </h2>
            <p className="text-sm text-muted-foreground mt-1">
              Control how you receive notifications
            </p>
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="notifications">Push Notifications</Label>
                <p className="text-sm text-muted-foreground">
                  Receive notifications in your browser
                </p>
              </div>
              <Switch
                id="notifications"
                checked={localSettings?.notifications_enabled || false}
                onCheckedChange={(checked) => handleSettingsChange('notifications_enabled', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="email-notifications">Email Notifications</Label>
                <p className="text-sm text-muted-foreground">
                  Receive notifications via email
                </p>
              </div>
              <Switch
                id="email-notifications"
                checked={localSettings?.email_notifications || false}
                onCheckedChange={(checked) => handleSettingsChange('email_notifications', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="sound">Sound Effects</Label>
                <p className="text-sm text-muted-foreground">
                  Play sounds for notifications and actions
                </p>
              </div>
              <Switch
                id="sound"
                checked={localSettings?.sound_enabled || false}
                onCheckedChange={(checked) => handleSettingsChange('sound_enabled', checked)}
              />
            </div>
          </div>
        </Card>

        {/* Auto-save Settings */}
        <Card className="p-6">
          <div className="mb-6">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Zap className="h-5 w-5" />
              Auto-save
            </h2>
            <p className="text-sm text-muted-foreground mt-1">
              Configure automatic saving of your work
            </p>
          </div>

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="auto-save">Enable Auto-save</Label>
                <p className="text-sm text-muted-foreground">
                  Automatically save your changes
                </p>
              </div>
              <Switch
                id="auto-save"
                checked={localSettings?.auto_save || false}
                onCheckedChange={(checked) => handleSettingsChange('auto_save', checked)}
              />
            </div>

            {localSettings?.auto_save && (
              <div>
                <Label htmlFor="auto-save-interval" className="flex items-center gap-2">
                  <Clock className="h-4 w-4" />
                  Auto-save Interval (seconds)
                </Label>
                <Input
                  id="auto-save-interval"
                  type="number"
                  min="10"
                  max="600"
                  value={localSettings?.auto_save_interval || 30}
                  onChange={(e) => handleSettingsChange('auto_save_interval', parseInt(e.target.value) || 30)}
                  className="mt-1"
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Minimum 10 seconds, maximum 10 minutes
                </p>
              </div>
            )}
          </div>
        </Card>
      </div>

      <Separator className="my-8" />

      <div className="flex justify-between items-center">
        <Button 
          variant="outline"
          onClick={handleResetToDefaults} 
          disabled={resetting || saving}
        >
          {resetting ? 'Resetting...' : 'Reset to Defaults'}
        </Button>
        
        <Button 
          onClick={saveSettings} 
          disabled={saving || resetting}
          size="lg"
        >
          {saving ? 'Saving...' : 'Save Settings'}
        </Button>
      </div>
    </div>
  );
};

export default Settings;