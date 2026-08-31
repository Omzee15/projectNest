import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Navbar } from '@/components/Navbar';
import { CreateProjectDialog } from '@/components/CreateProjectDialog';
import { ArrowRight, Users, Calendar, Zap, Download } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { isUserLoggedIn } from '@/utils/auth';

const Index = () => {
  const navigate = useNavigate();

  const handleGetStarted = () => {
    if (isUserLoggedIn()) {
      navigate('/projects');
    } else {
      navigate('/login');
    }
  };

  const features = [
    {
      icon: Users,
      title: 'Team Collaboration',
      description: 'Work together seamlessly with your team members',
    },
    {
      icon: Calendar,
      title: 'Project Timeline',
      description: 'Track deadlines and manage project schedules',
    },
    {
      icon: Zap,
      title: 'Quick Actions',
      description: 'Create and organize tasks with intuitive drag & drop',
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      
      <main className="container mx-auto px-4 sm:px-6 py-8 sm:py-16">
        {/* Hero Section */}
        <div className="text-center max-w-4xl mx-auto mb-12 sm:mb-16">
          <h1 className="text-3xl sm:text-4xl md:text-5xl font-bold text-foreground mb-4 sm:mb-6">
            Organize Your Projects
            <span className="block text-muted-foreground mt-2">Like Never Before</span>
          </h1>
          <p className="text-base sm:text-lg md:text-xl text-muted-foreground mb-6 sm:mb-8 max-w-2xl mx-auto px-4">
            A clean, intuitive project management platform inspired by Trello. 
            Create lists, manage tasks, and collaborate with your team effortlessly.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 justify-center items-center">
            <Button
              size="lg"
              onClick={handleGetStarted}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 sm:px-8 w-full sm:w-auto"
            >
              Get Started <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              onClick={() => navigate('/download')}
              className="px-6 sm:px-8 w-full sm:w-auto"
            >
              <Download className="mr-2 h-4 w-4" />
              Get the Menu Bar App
            </Button>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-6 md:gap-8 mb-12 sm:mb-16">
          {features.map((feature, index) => (
            <Card key={index} className="text-center hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="mx-auto w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                  <feature.icon className="h-6 w-6 text-primary" />
                </div>
                <CardTitle className="text-xl">{feature.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  {feature.description}
                </CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* CTA Section */}
        <div className="text-center bg-muted rounded-lg py-12 sm:py-16 px-4 sm:px-8">
          <h2 className="text-2xl sm:text-3xl font-bold text-foreground mb-4">
            Ready to Get Started?
          </h2>
          <p className="text-base sm:text-lg text-muted-foreground mb-6 sm:mb-8 max-w-2xl mx-auto">
            Experience the power of organized project management. 
            Create your first project and start collaborating today.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" onClick={() => navigate('/projects')} className="w-full sm:w-auto">
              View All Projects
            </Button>
            <CreateProjectDialog 
              trigger={
                <Button size="lg" variant="outline" className="w-full sm:w-auto">
                  Create New Project
                </Button>
              }
            />
          </div>
        </div>
      </main>
    </div>
  );
};

export default Index;
