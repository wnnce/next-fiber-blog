import { createContext, Dispatch, SetStateAction } from 'react'
import { User } from '@/lib/types'

export interface CommentState {
  type: number;
  articleId?: number;
  topicId?: number;
  user?: User
}

export interface StateContextProps {
  state: CommentState;
  setState: Dispatch<SetStateAction<CommentState>>
}

export const StateContext = createContext<{
  state: CommentState,
  setState: Dispatch<SetStateAction<CommentState>>
}>({ state: { type: 0 }, setState: pre => pre });