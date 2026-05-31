export interface StreakStats {
  current: number
  longest: number
  lastDate: string
}

export interface DayStats {
  reviews: number
}

export interface TotalStats {
  reviews: number
  cardsStudied: number
}

export interface StudyStats {
  streak: StreakStats
  today: DayStats
  total: TotalStats
}
