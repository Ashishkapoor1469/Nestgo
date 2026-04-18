import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export func cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
