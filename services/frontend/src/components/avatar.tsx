import { useMemo } from 'react'

export const AvatarIcon = ({ name }: { name: string }) => {
  const initials = useMemo(() => {
    return name
      .split(/[-_\s]/)
      .map(word => word[0].toUpperCase())
      .join('')
      .slice(0, 2)
  }, [name])

  const color = useMemo(() => {
    const colors = [
      'bg-red-500',
      'bg-green-500',
      'bg-blue-500',
      'bg-yellow-500',
      'bg-pink-500',
      'bg-purple-500',
      'bg-orange-500',
      'bg-teal-500',
    ]

    const stupidHash = Array.from(name).reduce((acc, char) => acc + char.charCodeAt(0), 0)
    return colors[stupidHash % colors.length]
  }, [name])

  return (
    <p className={`w-full h-full text-white flex items-center rounded-md justify-center font-medium text-sm sm:text-base ${color}`}>
      {initials}
    </p>
  )
}
