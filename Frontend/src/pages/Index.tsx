import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Navbar } from '@/components/Navbar';
import { CreateProjectDialog } from '@/components/CreateProjectDialog';
import { ArrowRight, Users, Calendar, Zap } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const Index = () => {
  const navigate = useNavigate();

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
      
      <main className="container mx-auto px-6 py-16">
        {/* Hero Section */}
        <div className="text-center max-w-4xl mx-auto mb-16">
          <h1 className="text-5xl font-bold text-foreground mb-6">
            Organize Your Projects
            <span className="block text-muted-foreground">Like Never Before</span>
          </h1>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            A clean, intuitive project management platform inspired by Trello. 
            Create lists, manage tasks, and collaborate with your team effortlessly.
          </p>
                    <Button 
            size="lg" 
            onClick={() => navigate('/login')}
            className="bg-blue-600 hover:bg-blue-700 text-white px-8"
          >
            Get Started <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-8 mb-16">
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
        <div className="text-center bg-muted rounded-lg py-16 px-8">
          <h2 className="text-3xl font-bold text-foreground mb-4">
            Ready to Get Started?
          </h2>
          <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
            Experience the power of organized project management. 
            Create your first project and start collaborating today.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" onClick={() => navigate('/projects')}>
              View All Projects
            </Button>
            <CreateProjectDialog 
              trigger={
                <Button size="lg" variant="outline">
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
