import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'

export function NotFound() {
  return (
    <div className="min-h-[50vh] flex flex-col items-center justify-center text-center px-4 py-20">
      <h1 className="text-4xl font-bold text-white mb-4">404 - Page Not Found</h1>
      <p className="text-slate-400 mb-8">The page you are looking for does not exist.</p>
      <Link to="/">
        <Button className="bg-cyan-500 hover:bg-cyan-600 text-white">
          Go Home
        </Button>
      </Link>
    </div>
  )
}
