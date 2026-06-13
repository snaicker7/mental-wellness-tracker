import React, { useEffect, useMemo, useState } from 'react'

const defaultApiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
const defaultUserId = import.meta.env.VITE_DEFAULT_USER_ID || 'Surya'

const moodFaces = ['🙂', '😌', '😐', '😟', '😩', '😵']
const moodLabels = {
  en: ['Calm', 'Flat', 'Muted', 'Worried', 'Stressed', 'Overloaded'],
  hi: ['शांत', 'सामान्य', 'मौन', 'चिंतित', 'तनावग्रस्त', 'अतिभारित']
}

const moodColors = ['#a5f39b', '#a3d9ff', '#d7ccff', '#ffe8d6', '#ff9470', '#e04c3f']

const translations = {
  en: {
    welcomeBack: "Welcome back",
    howAreYou: "How are you feeling today?",
    dashboardSubtitle: "A bright, mobile-first dashboard for tracking mood, stress, and study balance in one place.",
    moodDashboard: "Mood dashboard",
    currentMood: "Current mood",
    entriesTracked: "Entries tracked",
    stressLevel: "Stress level",
    trendingHigher: "Trending higher",
    trendingLower: "Trending lower",
    lastEntry: "Last",
    noEntriesYet: "No entries yet",
    weeklyMoodPulse: "Weekly mood pulse",
    wellnessTrend: "Wellness trend",
    chartSubtitle: "Mood rises in blue, stress rises in coral. Hover or tab to view data details.",
    chartEmptyMsg: "Add a few entries to see your mood and stress trend line appear here.",
    moodLegend: "Mood",
    stressLegend: "Stress",
    moodCalendar: "Mood calendar",
    thisMonth: "This month",
    monthlyView: "Monthly view",
    noCalendarData: "No calendar data yet.",
    todaysSnapshot: "Today's snapshot",
    quickInsight: "Quick insight",
    steadyMoodMsg: "Your current mood is {moodLabel}, with stress at {stress} and mood at {mood}.",
    noInsightMsg: "Add a few journal entries and the app will start shaping personalized insights here.",
    energy: "Energy",
    recentEntries: "Recent entries",
    savedNotes: "Saved notes",
    noNotesYet: "No saved entries yet.",
    backendUrl: "Backend URL",
    userId: "User ID",
    moodLabel: "Mood",
    stressLabel: "Stress",
    comfort: "Comfort",
    breathe: "Breathe",
    grounding: "Grounding",
    comfortKicker: "Empathetic Companion",
    comfortTitle: "Daily Companion Guide",
    crisisKicker: "Immediate Support",
    crisisTitle: "Wellness Help Resource",
    crisisWarning: "We are here for you.",
    aasra: "AASRA Helpline",
    vandrevala: "Vandrevala Foundation",
    crisisSubtext: "We recommend pausing and reaching out to a professional mental health counselor, doctor, or a trusted loved one.",
    gotIt: "Got it",
    journalInsights: "Journal Insights",
    stressTriggerAnalysis: "Stress Trigger Analysis",
    extractedTriggers: "Extracted Triggers",
    wellnessSummary: "Wellness Summary",
    empatheticGuidance: "Empathetic Guidance",
    scaleLabel: "Scale of 1 - 10",
    sleep: "Sleep",
    study: "Study",
    countdown: "Countdown",
    noTriggers: "No significant stress triggers detected in this entry.",
    breathingPractice: "Breathing Practice",
    mindfulGrounding: "Mindful Grounding",
    analyzingProgress: "Analyzing stress indicators & patterns...",
    stressScore: "Stress Score",
    openMenu: "Open menu",
    closeModal: "Close modal",
    writeJournal: "Write Today's Journal",
    journalPlaceholder: "How was your day? Write down your thoughts, achievements, or worries...",
    moodLevel: "Mood Level (1-10)",
    energyLevel: "Energy Level (1-10)",
    sleepDuration: "Sleep Duration (hours)",
    studyHours: "Study Hours (hours)",
    examCountdown: "Exam Countdown (days)",
    submitEntry: "Submit Entry",
    submitting: "Submitting...",
    entrySuccess: "Journal entry saved successfully!",
    entryError: "Failed to save entry: {error}",
    quickMoodSuccess: "Logged quick mood: {mood}!",
    validationError: "Please write at least a few words in your journal.",
    // New translations for dashboard redesign
    landingTitle: "Not Sure About Your Mood?",
    landingSubtitle: "Let us help you track, analyze, and build healthy coping habits today.",
    letUsHelp: "Let Us Help! â†’",
    notSureMood: "Not Sure About Your Mood?",
    sleepTitle: "Sleep Duration",
    stressTitle: "Stress Indicator",
    quizTitle: "Quick Wellness Quiz",
    questionNum: "Question {num} of {total}",
    yes: "Yes",
    no: "No",
    quizReset: "Restart Quiz",
    quizCompl: "Well done!",
    quizResult0: "You seem to be doing great! Keep it up.",
    quizResult1: "You are holding up well, but don't forget to take breaks.",
    quizResult2: "You might be feeling a bit overwhelmed. Try to get more sleep.",
    quizResult3: "Your stress levels seem high. Consider checking out the Coping Companion guide or talking to someone.",
    monthlySummary: "Monthly Mood Summary",
    activity: "Activity",
    steps: "steps",
    therapy: "Therapy",
    sessions: "sessions",
    discipline: "Discipline",
    focusScore: "focus score",
    tabHome: "Home",
    tabCalendar: "Calendar",
    tabTrends: "Trends",
    tabCoping: "Coping"
  },
  hi: {
    welcomeBack: "à¤†à¤ªà¤•à¤¾ à¤¸à¥�à¤µà¤¾à¤—à¤¤ à¤¹à¥ˆ",
    howAreYou: "à¤†à¤œ à¤†à¤ª à¤•à¥ˆà¤¸à¤¾ à¤®à¤¹à¤¸à¥‚à¤¸ à¤•à¤° à¤°à¤¹à¥‡ à¤¹à¥ˆà¤‚?",
    dashboardSubtitle: "à¤®à¥‚à¤¡, à¤¤à¤¨à¤¾à¤µ à¤”à¤° à¤…à¤§à¥�à¤¯à¤¯à¤¨ à¤¸à¤‚à¤¤à¥�à¤²à¤¨ à¤•à¥‹ à¤Ÿà¥�à¤°à¥ˆà¤• à¤•à¤°à¤¨à¥‡ à¤•à¥‡ à¤²à¤¿à¤� à¤�à¤• à¤‰à¤œà¥�à¤œà¥�à¤µà¤², à¤®à¥‹à¤¬à¤¾à¤‡à¤²-à¤«à¤°à¥�à¤¸à¥�à¤Ÿ à¤¡à¥ˆà¤¶à¤¬à¥‹à¤°à¥�à¤¡à¥¤",
    moodDashboard: "à¤®à¥‚à¤¡ à¤¡à¥ˆà¤¶à¤¬à¥‹à¤°à¥�à¤¡",
    currentMood: "à¤µà¤°à¥�à¤¤à¤®à¤¾à¤¨ à¤®à¥‚à¤¡",
    entriesTracked: "à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿à¤¯à¤¾à¤‚ à¤Ÿà¥�à¤°à¥ˆà¤• à¤•à¥€ à¤—à¤ˆà¤‚",
    stressLevel: "à¤¤à¤¨à¤¾à¤µ à¤•à¤¾ à¤¸à¥�à¤¤à¤°",
    trendingHigher: "à¤¬à¤¢à¤¼ à¤°à¤¹à¤¾ à¤¹à¥ˆ",
    trendingLower: "à¤•à¤® à¤¹à¥‹ à¤°à¤¹à¤¾ à¤¹à¥ˆ",
    lastEntry: "à¤…à¤‚à¤¤à¤¿à¤®",
    noEntriesYet: "à¤•à¥‹à¤ˆ à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤¨à¤¹à¥€à¤‚",
    weeklyMoodPulse: "à¤¸à¤¾à¤ªà¥�à¤¤à¤¾à¤¹à¤¿à¤• à¤®à¥‚à¤¡ à¤ªà¤²à¥�à¤¸",
    wellnessTrend: "à¤•à¤²à¥�à¤¯à¤¾à¤£ à¤•à¤¾ à¤°à¥�à¤�à¤¾à¤¨",
    chartSubtitle: "à¤¨à¥€à¤²à¥‡ à¤°à¤‚à¤— à¤®à¥‡à¤‚ à¤®à¥‚à¤¡ à¤¬à¤¢à¤¼à¤¤à¤¾ à¤¹à¥ˆ, à¤•à¥‹à¤°à¤² à¤®à¥‡à¤‚ à¤¤à¤¨à¤¾à¤µ à¤¬à¤¢à¤¼à¤¤à¤¾ à¤¹à¥ˆà¥¤ à¤µà¤¿à¤µà¤°à¤£ à¤¦à¥‡à¤–à¤¨à¥‡ à¤•à¥‡ à¤²à¤¿à¤� à¤¹à¥‹à¤µà¤° à¤•à¤°à¥‡à¤‚ à¤¯à¤¾ à¤Ÿà¥ˆà¤¬ à¤•à¤°à¥‡à¤‚à¥¤",
    chartEmptyMsg: "à¤…à¤ªà¤¨à¥€ à¤®à¥‚à¤¡ à¤”à¤° à¤¤à¤¨à¤¾à¤µ à¤•à¥€ à¤Ÿà¥�à¤°à¥‡à¤‚à¤¡ à¤²à¤¾à¤‡à¤¨ à¤¦à¥‡à¤–à¤¨à¥‡ à¤•à¥‡ à¤²à¤¿à¤� à¤•à¥�à¤› à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿à¤¯à¤¾à¤‚ à¤œà¥‹à¤¡à¤¼à¥‡à¤‚à¥¤",
    moodLegend: "à¤®à¥‚à¤¡",
    stressLegend: "à¤¤à¤¨à¤¾à¤µ",
    moodCalendar: "à¤®à¥‚à¤¡ à¤•à¥ˆà¤²à¥‡à¤‚à¤¡à¤°",
    thisMonth: "à¤‡à¤¸ à¤®à¤¹à¥€à¤¨à¥‡",
    monthlyView: "à¤®à¤¾à¤¸à¤¿à¤• à¤¦à¥ƒà¤¶à¥�à¤¯",
    noCalendarData: "à¤…à¤­à¥€ à¤•à¥‹à¤ˆ à¤•à¥ˆà¤²à¥‡à¤‚à¤¡à¤° à¤¡à¥‡à¤Ÿà¤¾ à¤¨à¤¹à¥€à¤‚ à¤¹à¥ˆà¥¤",
    todaysSnapshot: "à¤†à¤œ à¤•à¤¾ à¤¸à¥�à¤¨à¥ˆà¤ªà¤¶à¥‰à¤Ÿ",
    quickInsight: "à¤¤à¥�à¤µà¤°à¤¿à¤¤ à¤…à¤‚à¤¤à¤°à¥�à¤¦à¥ƒà¤·à¥�à¤Ÿà¤¿",
    steadyMoodMsg: "à¤†à¤ªà¤•à¤¾ à¤µà¤°à¥�à¤¤à¤®à¤¾à¤¨ à¤®à¥‚à¤¡ {moodLabel} à¤¹à¥ˆ, à¤œà¤¿à¤¸à¤®à¥‡à¤‚ à¤¤à¤¨à¤¾à¤µ {stress} à¤ªà¤° à¤”à¤° à¤®à¥‚à¤¡ {mood} à¤ªà¤° à¤¹à¥ˆà¥¤",
    noInsightMsg: "à¤•à¥�à¤› à¤œà¤°à¥�à¤¨à¤² à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿à¤¯à¤¾à¤‚ à¤œà¥‹à¤¡à¤¼à¥‡à¤‚ à¤”à¤° à¤�à¤ª à¤¯à¤¹à¤¾à¤‚ à¤µà¥�à¤¯à¤•à¥�à¤¤à¤¿à¤—à¤¤ à¤…à¤‚à¤¤à¤°à¥�à¤¦à¥ƒà¤·à¥�à¤Ÿà¤¿ à¤¦à¤¿à¤–à¤¾à¤¨à¤¾ à¤¶à¥�à¤°à¥‚ à¤•à¤° à¤¦à¥‡à¤—à¤¾à¥¤",
    energy: "à¤Šà¤°à¥�à¤œà¤¾",
    recentEntries: "à¤¹à¤¾à¤²à¤¿à¤¯à¤¾ à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿à¤¯à¤¾à¤‚",
    savedNotes: "à¤¸à¤¹à¥‡à¤œà¥‡ à¤—à¤� à¤¨à¥‹à¤Ÿ",
    noNotesYet: "à¤…à¤­à¥€ à¤•à¥‹à¤ˆ à¤¸à¤¹à¥‡à¤œà¥€ à¤—à¤ˆ à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤¨à¤¹à¥€à¤‚ à¤¹à¥ˆà¥¤",
    backendUrl: "à¤¬à¥ˆà¤•à¤�à¤‚à¤¡ URL",
    userId: "à¤¯à¥‚à¤œà¤¼à¤° ID",
    moodLabel: "à¤®à¥‚à¤¡",
    stressLabel: "à¤¤à¤¨à¤¾à¤µ",
    comfort: "à¤¸à¤¾à¤‚à¤¤à¥�à¤µà¤¨à¤¾",
    breathe: "à¤¸à¤¾à¤‚à¤¸",
    grounding: "à¤—à¥�à¤°à¤¾à¤‰à¤‚à¤¡à¤¿à¤‚à¤—",
    comfortKicker: "à¤¸à¤¹à¤¾à¤¨à¥�à¤­à¥‚à¤¤à¤¿à¤ªà¥‚à¤°à¥�à¤£ à¤¸à¤¾à¤¥à¥€",
    comfortTitle: "à¤¦à¥ˆà¤¨à¤¿à¤• à¤¸à¤¾à¤¥à¥€ à¤—à¤¾à¤‡à¤¡",
    crisisKicker: "à¤¤à¤¤à¥�à¤•à¤¾à¤² à¤¸à¤¹à¤¾à¤¯à¤¤à¤¾",
    crisisTitle: "à¤•à¤²à¥�à¤¯à¤¾à¤£ à¤¸à¤¹à¤¾à¤¯à¤¤à¤¾ à¤¸à¤‚à¤¸à¤¾à¤§à¤¨",
    crisisWarning: "à¤¹à¤® à¤†à¤ªà¤•à¥‡ à¤²à¤¿à¤� à¤¯à¤¹à¤¾à¤� à¤¹à¥ˆà¤‚à¥¤",
    aasra: "à¤†à¤¸à¤°à¤¾ à¤¹à¥‡à¤²à¥�à¤ªà¤²à¤¾à¤‡à¤¨",
    vandrevala: "à¤µà¤¾à¤‚à¤¡à¥�à¤°à¥‡à¤µà¤¾à¤²à¤¾ à¤«à¤¾à¤‰à¤‚à¤¡à¥‡à¤¶à¤¨",
    crisisSubtext: "à¤¹à¤® à¤•à¥�à¤› à¤¸à¤®à¤¯ à¤°à¥�à¤•à¤¨à¥‡ à¤”à¤° à¤•à¤¿à¤¸à¥€ à¤ªà¥‡à¤¶à¥‡à¤µà¤° à¤®à¤¾à¤¨à¤¸à¤¿à¤• à¤¸à¥�à¤µà¤¾à¤¸à¥�à¤¥à¥�à¤¯ à¤¸à¤²à¤¾à¤¹à¤•à¤¾à¤°, à¤¡à¥‰à¤•à¥�à¤Ÿà¤° à¤¯à¤¾ à¤•à¤¿à¤¸à¥€ à¤µà¤¿à¤¶à¥�à¤µà¤¸à¤¨à¥€à¤¯ à¤ªà¥�à¤°à¤¿à¤¯à¤œà¤¨ à¤¸à¥‡ à¤¸à¤‚à¤ªà¤°à¥�à¤• à¤•à¤°à¤¨à¥‡ à¤•à¥€ à¤¸à¤²à¤¾à¤¹ à¤¦à¥‡à¤¤à¥‡ à¤¹à¥ˆà¤‚à¥¤",
    gotIt: "à¤¸à¤®à¤� à¤—à¤�",
    journalInsights: "à¤œà¤°à¥�à¤¨à¤² à¤…à¤‚à¤¤à¤°à¥�à¤¦à¥ƒà¤·à¥�à¤Ÿà¤¿",
    stressTriggerAnalysis: "à¤¤à¤¨à¤¾à¤µ à¤Ÿà¥�à¤°à¤¿à¤—à¤° à¤µà¤¿à¤¶à¥�à¤²à¥‡à¤·à¤£",
    extractedTriggers: "à¤¨à¤¿à¤•à¤¾à¤²à¥‡ à¤—à¤� à¤Ÿà¥�à¤°à¤¿à¤—à¤°à¥�à¤¸",
    wellnessSummary: "à¤•à¤²à¥�à¤¯à¤¾à¤£ à¤¸à¤¾à¤°à¤¾à¤‚à¤¶",
    empatheticGuidance: "à¤¸à¤¹à¤¾à¤¨à¥�à¤­à¥‚à¤¤à¤¿à¤ªà¥‚à¤°à¥�à¤£ à¤®à¤¾à¤°à¥�à¤—à¤¦à¤°à¥�à¤¶à¤¨",
    scaleLabel: "1 à¤¸à¥‡ 10 à¤•à¤¾ à¤ªà¥ˆà¤®à¤¾à¤¨à¤¾",
    sleep: "à¤¨à¥€à¤‚à¤¦",
    study: "à¤…à¤§à¥�à¤¯à¤¯à¤¨",
    countdown: "à¤•à¤¾à¤‰à¤‚à¤Ÿà¤¡à¤¾à¤‰à¤¨",
    noTriggers: "à¤‡à¤¸ à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤®à¥‡à¤‚ à¤•à¥‹à¤ˆ à¤®à¤¹à¤¤à¥�à¤µà¤ªà¥‚à¤°à¥�à¤£ à¤¤à¤¨à¤¾à¤µ à¤Ÿà¥�à¤°à¤¿à¤—à¤° à¤¨à¤¹à¥€à¤‚ à¤®à¤¿à¤²à¤¾à¥¤",
    breathingPractice: "à¤¶à¥�à¤µà¤¸à¤¨ à¤…à¤­à¥�à¤¯à¤¾à¤¸",
    mindfulGrounding: "à¤®à¤¾à¤‡à¤‚à¤¡à¤«à¥�à¤² à¤—à¥�à¤°à¤¾à¤‰à¤‚à¤¡à¤¿à¤‚à¤—",
    analyzingProgress: "à¤¤à¤¨à¤¾à¤µ à¤¸à¤‚à¤•à¥‡à¤¤à¤•à¥‹à¤‚ à¤”à¤° à¤ªà¥ˆà¤Ÿà¤°à¥�à¤¨ à¤•à¤¾ à¤µà¤¿à¤¶à¥�à¤²à¥‡à¤·à¤£ à¤•à¤¿à¤¯à¤¾ à¤œà¤¾ à¤°à¤¹à¤¾ à¤¹à¥ˆ...",
    stressScore: "à¤¤à¤¨à¤¾à¤µ à¤¸à¥�à¤•à¥‹à¤°",
    openMenu: "à¤®à¥‡à¤¨à¥�à¤¯à¥‚ à¤–à¥‹à¤²à¥‡à¤‚",
    closeModal: "à¤¬à¤‚à¤¦ à¤•à¤°à¥‡à¤‚",
    writeJournal: "à¤†à¤œ à¤•à¥€ à¤œà¤°à¥�à¤¨à¤² à¤²à¤¿à¤–à¥‡à¤‚",
    journalPlaceholder: "à¤†à¤ªà¤•à¤¾ à¤¦à¤¿à¤¨ à¤•à¥ˆà¤¸à¤¾ à¤°à¤¹à¤¾? à¤…à¤ªà¤¨à¥‡ à¤µà¤¿à¤šà¤¾à¤°, à¤‰à¤ªà¤²à¤¬à¥�à¤§à¤¿à¤¯à¤¾à¤‚ à¤¯à¤¾ à¤šà¤¿à¤‚à¤¤à¤¾à¤�à¤‚ à¤²à¤¿à¤–à¥‡à¤‚...",
    moodLevel: "à¤®à¥‚à¤¡ à¤¸à¥�à¤¤à¤° (1-10)",
    energyLevel: "à¤Šà¤°à¥�à¤œà¤¾ à¤¸à¥�à¤¤à¤° (1-10)",
    sleepDuration: "à¤¨à¥€à¤‚à¤¦ à¤•à¥€ à¤…à¤µà¤§à¤¿ (à¤˜à¤‚à¤Ÿà¥‡)",
    studyHours: "à¤…à¤§à¥�à¤¯à¤¯à¤¨ à¤•à¥‡ à¤˜à¤‚à¤Ÿà¥‡",
    examCountdown: "à¤ªà¤°à¥€à¤•à¥�à¤·à¤¾ à¤•à¤¾à¤‰à¤‚à¤Ÿà¤¡à¤¾à¤‰à¤¨ (à¤¦à¤¿à¤¨)",
    submitEntry: "à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤œà¤®à¤¾ à¤•à¤°à¥‡à¤‚",
    submitting: "à¤œà¤®à¤¾ à¤¹à¥‹ à¤°à¤¹à¤¾ à¤¹à¥ˆ...",
    entrySuccess: "à¤œà¤°à¥�à¤¨à¤² à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤¸à¤«à¤²à¤¤à¤¾à¤ªà¥‚à¤°à¥�à¤µà¤• à¤¸à¤¹à¥‡à¤œà¥€ à¤—à¤ˆ!",
    entryError: "à¤ªà¥�à¤°à¤µà¤¿à¤·à¥�à¤Ÿà¤¿ à¤¸à¤¹à¥‡à¤œà¤¨à¥‡ à¤®à¥‡à¤‚ à¤µà¤¿à¤«à¤²: {error}",
    quickMoodSuccess: "à¤¤à¥�à¤µà¤°à¤¿à¤¤ à¤®à¥‚à¤¡ à¤¦à¤°à¥�à¤œ à¤•à¤¿à¤¯à¤¾ à¤—à¤¯à¤¾: {mood}!",
    validationError: "à¤•à¥ƒà¤ªà¤¯à¤¾ à¤…à¤ªà¤¨à¥€ à¤œà¤°à¥�à¤¨à¤² à¤®à¥‡à¤‚ à¤•à¤® à¤¸à¥‡ à¤•à¤® à¤•à¥�à¤› à¤¶à¤¬à¥�à¤¦ à¤²à¤¿à¤–à¥‡à¤‚à¥¤",
    // New translations for dashboard redesign
    landingTitle: "à¤…à¤ªà¤¨à¥‡ à¤®à¥‚à¤¡ à¤•à¥‹ à¤²à¥‡à¤•à¤° à¤…à¤¨à¤¿à¤¶à¥�à¤šà¤¿à¤¤ à¤¹à¥ˆà¤‚?",
    landingSubtitle: "à¤†à¤‡à¤� à¤¹à¤® à¤†à¤œ à¤†à¤ªà¤•à¥‹ à¤Ÿà¥�à¤°à¥ˆà¤• à¤•à¤°à¤¨à¥‡, à¤µà¤¿à¤¶à¥�à¤²à¥‡à¤·à¤£ à¤•à¤°à¤¨à¥‡ à¤”à¤° à¤¸à¥�à¤µà¤¸à¥�à¤¥ à¤®à¥�à¤•à¤¾à¤¬à¤²à¤¾ à¤•à¤°à¤¨à¥‡ à¤•à¥€ à¤†à¤¦à¤¤à¥‡à¤‚ à¤¬à¤¨à¤¾à¤¨à¥‡ à¤®à¥‡à¤‚ à¤®à¤¦à¤¦ à¤•à¤°à¥‡à¤‚à¥¤",
    letUsHelp: "à¤†à¤‡à¤� à¤¹à¤® à¤®à¤¦à¤¦ à¤•à¤°à¥‡à¤‚! â†’",
    notSureMood: "à¤…à¤ªà¤¨à¥‡ à¤®à¥‚à¤¡ à¤•à¥‹ à¤²à¥‡à¤•à¤° à¤…à¤¨à¤¿à¤¶à¥�à¤šà¤¿à¤¤ à¤¹à¥ˆà¤‚?",
    sleepTitle: "à¤¨à¥€à¤‚à¤¦ à¤•à¥€ à¤…à¤µà¤§à¤¿",
    stressTitle: "à¤¤à¤¨à¤¾à¤µ à¤¸à¥‚à¤šà¤•",
    quizTitle: "à¤¤à¥�à¤µà¤°à¤¿à¤¤ à¤•à¤²à¥�à¤¯à¤¾à¤£ à¤ªà¥�à¤°à¤¶à¥�à¤¨à¥‹à¤¤à¥�à¤¤à¤°à¥€",
    questionNum: "à¤ªà¥�à¤°à¤¶à¥�à¤¨ {num} à¤•à¤¾ {total}",
    yes: "à¤¹à¤¾à¤�",
    no: "à¤¨à¤¹à¥€à¤‚",
    quizReset: "à¤ªà¥�à¤°à¤¶à¥�à¤¨à¥‹à¤¤à¥�à¤¤à¤°à¥€ à¤ªà¥�à¤¨à¤°à¤¾à¤°à¤‚à¤­ à¤•à¤°à¥‡à¤‚",
    quizCompl: "à¤¬à¤¹à¥�à¤¤ à¤¬à¤¢à¤¼à¤¿à¤¯à¤¾!",
    quizResult0: "à¤†à¤ª à¤¬à¤¹à¥�à¤¤ à¤…à¤šà¥�à¤›à¤¾ à¤•à¤° à¤°à¤¹à¥‡ à¤¹à¥ˆà¤‚! à¤‡à¤¸à¥‡ à¤œà¤¾à¤°à¥€ à¤°à¤–à¥‡à¤‚à¥¤",
    quizResult1: "à¤†à¤ª à¤…à¤šà¥�à¤›à¤¾ à¤•à¤° à¤°à¤¹à¥‡ à¤¹à¥ˆà¤‚, à¤²à¥‡à¤•à¤¿à¤¨ à¤¬à¥�à¤°à¥‡à¤• à¤²à¥‡à¤¨à¤¾ à¤¨ à¤­à¥‚à¤²à¥‡à¤‚à¥¤",
    quizResult2: "à¤†à¤ª à¤¥à¥‹à¤¡à¤¼à¤¾ à¤…à¤­à¤¿à¤­à¥‚à¤¤ à¤®à¤¹à¤¸à¥‚à¤¸ à¤•à¤° à¤°à¤¹à¥‡ à¤¹à¥‹à¤‚à¤—à¥‡à¥¤ à¤…à¤§à¤¿à¤• à¤¸à¥‹à¤¨à¥‡ à¤•à¥€ à¤•à¥‹à¤¶à¤¿à¤¶ à¤•à¤°à¥‡à¤‚à¥¤",
    quizResult3: "à¤†à¤ªà¤•à¤¾ à¤¤à¤¨à¤¾à¤µ à¤•à¤¾ à¤¸à¥�à¤¤à¤° à¤•à¤¾à¤«à¥€ à¤…à¤§à¤¿à¤• à¤²à¤— à¤°à¤¹à¤¾ à¤¹à¥ˆà¥¤ à¤¸à¤¾à¤‚à¤¤à¥�à¤µà¤¨à¤¾ à¤—à¤¾à¤‡à¤¡ à¤¦à¥‡à¤–à¥‡à¤‚ à¤¯à¤¾ à¤•à¤¿à¤¸à¥€ à¤¸à¥‡ à¤¬à¤¾à¤¤ à¤•à¤°à¥‡à¤‚à¥¤",
    monthlySummary: "à¤®à¤¾à¤¸à¤¿à¤• à¤®à¥‚à¤¡ à¤¸à¤¾à¤°à¤¾à¤‚à¤¶",
    activity: "à¤—à¤¤à¤¿à¤µà¤¿à¤§à¤¿",
    steps: "à¤•à¤¦à¤®",
    therapy: "à¤¥à¥‡à¤°à¥‡à¤ªà¥€",
    sessions: "à¤¸à¤¤à¥�à¤°",
    discipline: "à¤…à¤¨à¥�à¤¶à¤¾à¤¸à¤¨",
    focusScore: "à¤«à¥‹à¤•à¤¸ à¤¸à¥�à¤•à¥‹à¤°",
    tabHome: "à¤¹à¥‹à¤®",
    tabCalendar: "à¤•à¥ˆà¤²à¥‡à¤‚à¤¡à¤°",
    tabTrends: "à¤°à¥�à¤�à¤¾à¤¨",
    tabCoping: "à¤¸à¤¾à¤‚à¤¤à¥�à¤µà¤¨à¤¾"
  }
}

function formatValue(value) {
  return Number.isFinite(value) ? value.toFixed(1) : '0.0'
}

function formatShortDate(dateString) {
  if (!dateString) return ''
  const parsed = new Date(`${dateString}T00:00:00`)
  if (Number.isNaN(parsed.getTime())) return dateString
  return parsed.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function getMoodState(moodAverage, lang = 'en') {
  const bucket = Math.max(0, Math.min(moodFaces.length - 1, Math.round(10 - moodAverage)))
  return {
    face: moodFaces[bucket],
    label: moodLabels[lang][bucket],
    color: moodColors[bucket],
  }
}

function buildPolyline(points, width, height, topPad, bottomPad, valueSelector) {
  if (points.length === 0) return ''

  const values = points.map(valueSelector)
  const min = Math.min(...values)
  const max = Math.max(...values)
  const span = max - min || 1
  const usableHeight = height - topPad - bottomPad
  const stepX = points.length > 1 ? width / (points.length - 1) : width

  return points
    .map((point, index) => {
      const x = points.length > 1 ? index * stepX : width / 2
      const normalized = (valueSelector(point) - min) / span
      const y = height - bottomPad - normalized * usableHeight
      return `${x},${y}`
    })
    .join(' ')
}

function fetchJson(url, options) {
  return fetch(url, options).then(async (response) => {
    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}`)
    }
    return response.json()
  })
}

function TrendChart({ points, lang = 'en' }) {
  const [activeIndex, setActiveIndex] = useState(null)
  const width = 760
  const height = 320
  const topPad = 28
  const bottomPad = 40
  const t = translations[lang]

  const moodLine = useMemo(
    () => buildPolyline(points, width, height, topPad, bottomPad, (point) => point.mood_average),
    [points],
  )
  const stressLine = useMemo(
    () => buildPolyline(points, width, height, topPad, bottomPad, (point) => point.stress_average),
    [points],
  )

  const moodExtremes = useMemo(() => {
    if (points.length === 0) return { min: 0, max: 0 }
    const values = points.map((point) => point.mood_average)
    return { min: Math.min(...values), max: Math.max(...values) }
  }, [points])

  const stressExtremes = useMemo(() => {
    if (points.length === 0) return { min: 0, max: 0 }
    const values = points.map((point) => point.stress_average)
    return { min: Math.min(...values), max: Math.max(...values) }
  }, [points])

  if (points.length === 0) {
    return (
      <figure className="chart-card chart-card-empty">
        <figcaption>
          <h2>{t.weeklyMoodPulse}</h2>
          <p>{t.chartEmptyMsg}</p>
        </figcaption>
        <div className="empty-chart-illustration" aria-hidden="true">
          <span>â—Ž</span>
          <span>â—Œ</span>
          <span>â—�</span>
        </div>
      </figure>
    )
  }

  const activePoint = activeIndex !== null ? points[activeIndex] : null

  return (
    <figure className="chart-card" aria-labelledby="trend-chart-title" aria-describedby="trend-chart-help">
      <div className="section-heading section-heading-row">
        <div>
          <p className="section-kicker">{t.wellnessTrend}</p>
          <h2 id="trend-chart-title">{t.weeklyMoodPulse}</h2>
          <p id="trend-chart-help">{t.chartSubtitle}</p>
        </div>
        <div className="legend" aria-hidden="true">
          <span><i className="legend-mood" />{t.moodLegend}</span>
          <span><i className="legend-stress" />{t.stressLegend}</span>
        </div>
      </div>

      <div style={{ position: 'relative' }}>
        <div className="chart-scroll" tabIndex="0" aria-label="Scrollable trend chart">
          <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label="Mood and stress trend chart">
            <defs>
              <linearGradient id="moodFill" x1="0" x2="0" y1="0" y2="1">
                <stop offset="0%" stopColor="rgba(40, 98, 255, 0.24)" />
                <stop offset="100%" stopColor="rgba(40, 98, 255, 0.02)" />
              </linearGradient>
              <linearGradient id="stressFill" x1="0" x2="0" y1="0" y2="1">
                <stop offset="0%" stopColor="rgba(222, 76, 63, 0.18)" />
                <stop offset="100%" stopColor="rgba(222, 76, 63, 0.02)" />
              </linearGradient>
            </defs>
            <line x1="0" y1={height - bottomPad} x2={width} y2={height - bottomPad} className="axis-line" />
            {[0, 1, 2, 3].map((tick) => {
              const y = topPad + tick * ((height - topPad - bottomPad) / 3)
              return <line key={tick} x1="0" y1={y} x2={width} y2={y} className="grid-line" />
            })}

            <path d={`M 0 ${height - bottomPad} L ${moodLine} L ${width} ${height - bottomPad} Z`} fill="url(#moodFill)" opacity="0.8" />
            <path d={`M 0 ${height - bottomPad} L ${stressLine} L ${width} ${height - bottomPad} Z`} fill="url(#stressFill)" opacity="0.8" />
            <polyline points={moodLine} className="trend-line mood-line" />
            <polyline points={stressLine} className="trend-line stress-line" />

            {points.map((point, index) => {
              const x = points.length > 1 ? index * (width / (points.length - 1)) : width / 2
              const moodRange = moodExtremes.max - moodExtremes.min || 1
              const stressRange = stressExtremes.max - stressExtremes.min || 1
              const moodY = height - bottomPad - ((point.mood_average - moodExtremes.min) / moodRange) * (height - topPad - bottomPad)
              const stressY = height - bottomPad - ((point.stress_average - stressExtremes.min) / stressRange) * (height - topPad - bottomPad)

              const isHovered = activeIndex === index
              const stepX = points.length > 1 ? width / (points.length - 1) : width
              const hitWidth = points.length > 1 ? stepX : width

              return (
                <g 
                  key={`${point.date}-${index}`}
                  tabIndex="0"
                  role="button"
                  aria-label={`${t.lastEntry}: ${formatShortDate(point.date)}, ${t.moodLabel}: ${formatValue(point.mood_average)}, ${t.stressLabel}: ${formatValue(point.stress_average)}, ${t.energy}: ${formatValue(point.energy_average)}`}
                  onMouseEnter={() => setActiveIndex(index)}
                  onMouseLeave={() => setActiveIndex(null)}
                  onFocus={() => setActiveIndex(index)}
                  onBlur={() => setActiveIndex(null)}
                  className={`chart-point-group ${isHovered ? 'active' : ''}`}
                  style={{ outline: 'none' }}
                >
                  {isHovered && (
                    <line 
                      x1={x} 
                      y1={topPad} 
                      x2={x} 
                      y2={height - bottomPad} 
                      stroke="var(--accent)" 
                      strokeDasharray="4,4" 
                      strokeWidth="1.5" 
                    />
                  )}
                  
                  <rect 
                    x={x - hitWidth / 2} 
                    y={topPad} 
                    width={hitWidth} 
                    height={height - topPad - bottomPad} 
                    fill="transparent" 
                    style={{ cursor: 'pointer' }}
                  />

                  <circle cx={x} cy={moodY} r={isHovered ? "7" : "5"} className="point mood-point" />
                  <circle cx={x} cy={stressY} r={isHovered ? "7" : "5"} className="point stress-point" />
                  
                  <text x={x} y={height - 12} textAnchor="middle" className="chart-date">
                    {point.date.slice(5)}
                  </text>
                </g>
              )
            })}
          </svg>
        </div>

        {activePoint && (
          <div 
            className="chart-tooltip fade-in"
            style={{
              position: 'absolute',
              left: `${points.length > 1 ? (activeIndex * 100) / (points.length - 1) : 50}%`,
              transform: 'translate(-50%, -105%)',
              top: `${topPad}px`,
              pointerEvents: 'none',
              zIndex: 10,
            }}
          >
            <strong>{formatShortDate(activePoint.date)}</strong>
            <div className="tooltip-metrics">
              <div className="tooltip-metric mood-field">
                <span>{t.moodLabel}:</span>
                <strong>{formatValue(activePoint.mood_average)}</strong>
              </div>
              <div className="tooltip-metric stress-field">
                <span>{t.stressLabel}:</span>
                <strong>{formatValue(activePoint.stress_average)}</strong>
              </div>
              <div className="tooltip-metric energy-field">
                <span>{t.energy}:</span>
                <strong>{formatValue(activePoint.energy_average)}</strong>
              </div>
            </div>
          </div>
        )}
      </div>
    </figure>
  )
}

function MoodGrid({ points, lang = 'en' }) {
  const t = translations[lang]
  const cells = useMemo(() => {
    return points.slice(-18).map((point) => ({
      ...point,
      moodState: getMoodState(point.mood_average, lang),
    }))
  }, [points, lang])

  const monthLabel = points.length > 0 ? new Date(`${points.at(-1).date}T00:00:00`).toLocaleString(lang, { month: 'long', year: 'numeric' }) : t.thisMonth

  return (
    <section className="panel mood-grid-panel" aria-labelledby="mood-grid-title">
      <div className="section-heading section-heading-row">
        <div>
          <p className="section-kicker">{t.moodCalendar}</p>
          <h2 id="mood-grid-title">{t.moodCalendar}</h2>
          <p>{monthLabel}</p>
        </div>
        <button type="button" className="ghost-chip" onClick={() => {}} aria-label="Mood calendar controls">
          <span aria-hidden="true">â—Œ</span>
          <span>{t.monthlyView}</span>
        </button>
      </div>

      <div className="mood-grid" role="list" aria-label="Mood history grid">
        {cells.length === 0 ? (
          <div className="mood-grid-empty" role="status">{t.noCalendarData}</div>
        ) : (
          cells.map((cell) => (
            <article
              key={cell.date}
              role="listitem"
              className="mood-cell"
              style={{ '--cell-color': cell.moodState.color }}
              aria-label={`${cell.date} ${t.moodLabel} ${formatValue(cell.mood_average)} ${t.stressLabel} ${formatValue(cell.stress_average)}`}
            >
              <span className="mood-cell-date">{cell.date.slice(5)}</span>
              <span className="mood-cell-face" aria-hidden="true">{cell.moodState.face}</span>
              <span className="mood-cell-label">{cell.moodState.label}</span>
            </article>
          ))
        )}
      </div>
    </section>
  )
}

function App() {
  const [apiBaseUrl, setApiBaseUrl] = useState(defaultApiBaseUrl)
  const [userId, setUserId] = useState(defaultUserId)
  const [points, setPoints] = useState([])
  const [entries, setEntries] = useState([])
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')

  const [selectedEntry, setSelectedEntry] = useState(null)
  const [analysis, setAnalysis] = useState(null)
  const [analysisLoading, setAnalysisLoading] = useState(false)
  const [analysisError, setAnalysisError] = useState('')

  const [copingData, setCopingData] = useState(null)
  const [copingTab, setCopingTab] = useState('motivation')
  const [lang, setLang] = useState('en')
  const t = translations[lang]

  // Form states for journal logging
  const [entryText, setEntryText] = useState('')
  const [moodLevel, setMoodLevel] = useState(6)
  const [energyLevel, setEnergyLevel] = useState(6)
  const [sleepHours, setSleepHours] = useState(7)
  const [studyHours, setStudyHours] = useState(6)
  const [examCountdownDays, setExamCountdownDays] = useState(10)
  
  const [submitStatus, setSubmitStatus] = useState('idle') // idle, submitting, success, error
  const [toastMessage, setToastMessage] = useState('')
  const [refreshCounter, setRefreshCounter] = useState(0)

  // Redesign tab control and settings form toggler
  const [activeTab, setActiveTab] = useState('landing')
  const [showSettings, setShowSettings] = useState(false)

  // Quiz states
  const [quizStep, setQuizStep] = useState(0)
  const [quizAnswers, setQuizAnswers] = useState([])

  useEffect(() => {
    if (toastMessage) {
      const timer = setTimeout(() => {
        setToastMessage('')
        setSubmitStatus('idle')
      }, 4000)
      return () => clearTimeout(timer)
    }
  }, [toastMessage])

  const showToast = (message, status = 'success') => {
    setToastMessage(message)
    setSubmitStatus(status)
  }

  const handleQuickMood = async (moodName, moodVal, energyVal) => {
    try {
      setSubmitStatus('submitting')
      const trimmedBase = apiBaseUrl.replace(/\/$/, '')
      const payload = {
        user_id: userId,
        entry_text: lang === 'en' ? `Quick mood check-in: ${moodName}` : `त्वरित मूड चेक-इन: ${moodName}`,
        mood_level: moodVal,
        energy_level: energyVal,
        sleep_hours: 8.0,
        study_hours: 4.0,
        exam_countdown_days: 10
      }
      
      await fetchJson(`${trimmedBase}/v1/entries`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })
      
      const successMsg = t.quickMoodSuccess.replace('{mood}', moodName)
      showToast(successMsg, 'success')
      setRefreshCounter(prev => prev + 1)
    } catch (err) {
      showToast(t.entryError.replace('{error}', err.message), 'error')
    }
  }

  const handleFullSubmit = async (e) => {
    e.preventDefault()
    
    const trimmedText = entryText.trim()
    if (!trimmedText) {
      showToast(t.validationError, 'error')
      return
    }

    try {
      setSubmitStatus('submitting')
      const trimmedBase = apiBaseUrl.replace(/\/$/, '')
      const payload = {
        user_id: userId,
        entry_text: trimmedText,
        mood_level: moodLevel,
        energy_level: energyLevel,
        sleep_hours: sleepHours,
        study_hours: studyHours,
        exam_countdown_days: examCountdownDays
      }

      await fetchJson(`${trimmedBase}/v1/entries`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      showToast(t.entrySuccess, 'success')
      
      // Reset form fields
      setEntryText('')
      setMoodLevel(6)
      setEnergyLevel(6)
      setSleepHours(7)
      setStudyHours(6)
      setExamCountdownDays(10)
      
      setRefreshCounter(prev => prev + 1)
    } catch (err) {
      showToast(t.entryError.replace('{error}', err.message), 'error')
    }
  }

  const handleEntryClick = async (entry) => {
    setSelectedEntry(entry)
    setAnalysis(null)
    setAnalysisError('')
    setAnalysisLoading(true)
    try {
      const trimmedBase = apiBaseUrl.replace(/\/$/, '')
      const data = await fetchJson(
        `${trimmedBase}/v1/analysis?user_id=${encodeURIComponent(userId)}&entry_id=${encodeURIComponent(entry.id)}`,
      )
      setAnalysis(data)
    } catch (err) {
      setAnalysisError(err.message || 'Failed to load stress analysis')
    } finally {
      setAnalysisLoading(false)
    }
  }

  useEffect(() => {
    let isActive = true

    async function loadDashboard() {
      setStatus('loading')
      setError('')

      try {
        const trimmedBase = apiBaseUrl.replace(/\/$/, '')
        const [trendPayload, entryPayload, copingPayload] = await Promise.all([
          fetchJson(`${trimmedBase}/v1/trends?user_id=${encodeURIComponent(userId)}`),
          fetchJson(`${trimmedBase}/v1/entries?user_id=${encodeURIComponent(userId)}`),
          fetchJson(`${trimmedBase}/v1/coping?user_id=${encodeURIComponent(userId)}`).catch(() => null),
        ])

        if (!isActive) return

        setPoints(trendPayload.points || [])
        setEntries(entryPayload.entries || [])
        setCopingData(copingPayload || null)
        setStatus('ready')
      } catch (requestError) {
        if (!isActive) return
        setError(requestError.message || 'Unable to load dashboard')
        setStatus('error')
      }
    }

    loadDashboard()
    return () => {
      isActive = false
    }
  }, [apiBaseUrl, userId, refreshCounter])

  const latestPoint = points.at(-1)
  const latestEntry = entries.at(-1)
  const latestMoodState = latestPoint ? getMoodState(latestPoint.mood_average, lang) : null
  const trendDelta = points.length > 1 ? latestPoint.mood_average - points[points.length - 2].mood_average : 0
  const topMoodLabel = latestMoodState ? latestMoodState.label : (lang === 'hi' ? 'तैयार' : 'Ready to capture')
  
  const shortcutLabels = {
    en: ['Happy', 'Calm', 'Tired', 'Stressed'],
    hi: ['खुश', 'शांत', 'थका हुआ', 'तनावग्रस्त']
  }

  // Quiz questions engine
  const quizQuestions = {
    en: [
      "Have you been sleeping at least 7-8 hours recently?",
      "Are you studying more than 8 hours daily?",
      "Do you take regular breaks during study sessions?",
      "Do you feel overwhelmed by your upcoming exams?"
    ],
    hi: [
      "क्या आप हाल ही में कम से कम 7-8 घंटे सो रहे हैं?",
      "क्या आप रोजाना 8 घंटे से अधिक पढ़ाई कर रहे हैं?",
      "क्या आप पढ़ाई के दौरान नियमित अंतराल पर ब्रेक लेते हैं?",
      "क्या आप आगामी परीक्षाओं को लेकर तनाव महसूस कर रहे हैं?"
    ]
  }

  // Quiz helper functions
  const handleQuizAnswer = (isYes) => {
    const nextAnswers = [...quizAnswers, isYes]
    setQuizAnswers(nextAnswers)
    setQuizStep(prev => prev + 1)
  }

  const handleQuizReset = () => {
    setQuizStep(0)
    setQuizAnswers([])
  }

  const quizResultText = useMemo(() => {
    if (quizAnswers.length < 4) return ''
    let score = 0
    if (!quizAnswers[0]) score++ // No sleep -> stressed
    if (quizAnswers[1]) score++  // Studying > 8h -> stressed
    if (!quizAnswers[2]) score++ // No breaks -> stressed
    if (quizAnswers[3]) score++  // Overwhelmed -> stressed

    if (score === 0) return t.quizResult0
    if (score === 1) return t.quizResult1
    if (score === 2) return t.quizResult2
    return t.quizResult3
  }, [quizAnswers, lang, t])

  // June 2026 calendar days generation
  const calendarDays = useMemo(() => {
    const days = []
    // May 31 (Previous Month Filler)
    days.push({ date: '2026-05-31', dayNum: 31, isCurrentMonth: false })
    // June 1 to 30
    for (let d = 1; d <= 30; d++) {
      const dayStr = d < 10 ? `0${d}` : `${d}`
      days.push({ date: `2026-06-${dayStr}`, dayNum: d, isCurrentMonth: true })
    }
    // July 1 to 4 (Next Month Filler)
    for (let d = 1; d <= 4; d++) {
      days.push({ date: `2026-07-0${d}`, dayNum: d, isCurrentMonth: false })
    }
    return days
  }, [])

  // Dominant mood computation
  const dominantMood = useMemo(() => {
    const juneEntries = entries.filter(e => e.created_at && e.created_at.startsWith('2026-06'))
    if (juneEntries.length === 0) return null
    
    const counts = [0, 0, 0, 0, 0, 0]
    juneEntries.forEach(entry => {
      const moodVal = entry.mood_level || 5
      const bucket = Math.max(0, Math.min(moodFaces.length - 1, Math.round(10 - moodVal)))
      counts[bucket]++
    })
    
    let maxIdx = 0
    let maxVal = -1
    for (let i = 0; i < counts.length; i++) {
      if (counts[i] > maxVal) {
        maxVal = counts[i]
        maxIdx = i
      }
    }
    return getMoodState(10 - maxIdx, lang)
  }, [entries, lang])

  const dominantMoodDesc = useMemo(() => {
    if (!dominantMood) return lang === 'hi' ? 'कोई प्रविष्टि अभी तक नहीं।' : 'No entries logged yet.'
    const label = dominantMood.label
    
    if (label === 'Calm' || label === 'शांत') {
      return lang === 'hi' ? 'भारत में आपका यह महीना काफी शांत और संतुलित रहा। इसे बनाए रखें!' : 'You had a very calm and balanced month. Keep up the good work!'
    } else if (label === 'Flat' || label === 'सामान्य') {
      return lang === 'hi' ? 'आपका मूड ज्यादातर स्थिर और सामान्य था।' : 'Your mood was mostly stable and flat.'
    } else if (label === 'Muted' || label === 'मौन') {
      return lang === 'hi' ? 'आप कुछ हद तक मौन या शांत महसूस कर रहे थे।' : 'You felt somewhat muted or quiet.'
    } else if (label === 'Worried' || label === 'चिंतित') {
      return lang === 'hi' ? 'आपने कई दिन चिंतित या बेचैन महसूस करने में बिताए।' : 'You spent several days feeling worried or anxious.'
    } else if (label === 'Stressed' || label === 'तनावग्रस्त') {
      return lang === 'hi' ? 'आपका तनाव का स्तर ऊंचा था। सांस लेने के अभ्यास को याद रखें।' : 'Your stress level was high. Remember to practice breathing.'
    } else {
      return lang === 'hi' ? 'आपका यह महीना अतिभारित रहा। कृपया आराम को प्राथमिकता दें।' : 'You had an overloaded month. Please prioritize rest.'
    }
  }, [dominantMood, lang])

  // Get sleep summary & stress indicator data for last 6 entries
  const last6Entries = useMemo(() => {
    return entries.slice(-6)
  }, [entries])

  const sleepHoursDisplay = useMemo(() => {
    if (entries.length === 0) return '8h 0m'
    const latestSleep = entries.at(-1).sleep_hours || 0
    const hrs = Math.floor(latestSleep)
    const mins = Math.round((latestSleep - hrs) * 60)
    return `${hrs}h ${mins}m`
  }, [entries])

  const stressIndicatorLabel = useMemo(() => {
    if (entries.length === 0) return lang === 'hi' ? 'सामान्य' : 'Low'
    const latest = entries.at(-1)
    const stressVal = 11 - (latest.mood_level || 6)
    if (stressVal <= 3) return lang === 'hi' ? 'कम' : 'Low'
    if (stressVal <= 7) return lang === 'hi' ? 'मध्यम' : 'Medium'
    return lang === 'hi' ? 'उच्च' : 'High'
  }, [entries, lang])

  return (
    <main className="app-shell">
      {activeTab === 'landing' ? (
        <section className="landing-container">
          <div className="landing-shapes">
            <svg className="floating-shape shape-circle" viewBox="0 0 100 100" width="80" height="80">
              <circle cx="50" cy="50" r="45" fill="#a5f39b" />
              <circle cx="35" cy="45" r="4" fill="#121212" />
              <circle cx="65" cy="45" r="4" fill="#121212" />
              <path d="M 35 65 Q 50 75 65 65" stroke="#121212" strokeWidth="4" fill="transparent" strokeLinecap="round" />
            </svg>
            <svg className="floating-shape shape-cloud" viewBox="0 0 100 60" width="100" height="60">
              <path d="M 20 40 A 20 20 0 0 1 50 15 A 15 15 0 0 1 80 30 A 15 15 0 0 1 80 50 A 15 15 0 0 1 20 50 Z" fill="#a3d9ff" />
              <path d="M 33 37 Q 38 41 33 45 M 63 37 Q 58 41 63 45" stroke="#121212" strokeWidth="3" fill="transparent" strokeLinecap="round" />
              <text x="75" y="20" fontSize="12" fontWeight="bold" fill="#2862ff">Zzz</text>
            </svg>
            <svg className="floating-shape shape-star" viewBox="0 0 100 100" width="70" height="70">
              <polygon points="50,5 64,36 98,36 70,57 81,91 50,70 19,91 30,57 2,36 36,36" fill="#ffd54f" />
              <circle cx="40" cy="45" r="3" fill="#121212" />
              <circle cx="60" cy="45" r="3" fill="#121212" />
              <line x1="40" y1="58" x2="60" y2="58" stroke="#121212" strokeWidth="3" strokeLinecap="round" />
            </svg>
          </div>

          <div className="landing-lang-switcher">
            <button type="button" className={`lang-btn ${lang === 'en' ? 'active' : ''}`} onClick={() => setLang('en')}>EN</button>
            <button type="button" className={`lang-btn ${lang === 'hi' ? 'active' : ''}`} onClick={() => setLang('hi')}>हिन्दी</button>
          </div>

          <h1 className="landing-title">{t.landingTitle}</h1>
          <p className="landing-subtitle">{t.landingSubtitle}</p>
          <button type="button" className="landing-cta-btn" onClick={() => setActiveTab('home')}>
            {t.letUsHelp}
          </button>
        </section>
      ) : (
        <>
          <header className="app-header">
            <div className="header-left">
              <div className="avatar" aria-hidden="true">
                <span>SM</span>
              </div>
              <div className="header-greeting">
                <p>{t.welcomeBack}, <strong>{userId}</strong></p>
                <span className="header-date">Jun 13, 2026</span>
              </div>
            </div>

            <div className="header-right">
              <div className="lang-switcher">
                <button type="button" className={`lang-btn ${lang === 'en' ? 'active' : ''}`} onClick={() => setLang('en')}>EN</button>
                <button type="button" className={`lang-btn ${lang === 'hi' ? 'active' : ''}`} onClick={() => setLang('hi')}>हिन्दी</button>
              </div>
              <button type="button" className="settings-toggle-btn" onClick={() => setShowSettings(!showSettings)} aria-label={t.openMenu}>
                ⚙️
              </button>
            </div>
          </header>

          {showSettings && (
            <div className="settings-panel-overlay fade-in">
              <div className="settings-panel">
                <div className="settings-header">
                  <h3>App Settings</h3>
                  <button type="button" className="close-btn" onClick={() => setShowSettings(false)}>&times;</button>
                </div>
                <form className="controls" onSubmit={(e) => e.preventDefault()}>
                  <label>
                    {t.backendUrl}
                    <input value={apiBaseUrl} onChange={(e) => setApiBaseUrl(e.target.value)} placeholder="http://localhost:8080" />
                  </label>
                  <label>
                    {t.userId}
                    <input value={userId} onChange={(e) => setUserId(e.target.value)} placeholder="Surya" />
                  </label>
                </form>
              </div>
            </div>
          )}

          <div className="tab-pane-container">
            {activeTab === 'home' && (
              <div className="tab-view home-view fade-in">
                <section className="welcome-greeting-section">
                  <h2>{lang === 'hi' ? 'नमस्ते Surya!' : 'Hello Surya!'}</h2>
                  <p>{t.howAreYou}</p>
                </section>

                <section className="quick-mood-row" aria-label="Quick mood shortcuts">
                  {shortcutLabels[lang].map((label, index) => {
                    const moodValues = [9, 10, 5, 2]
                    const energyValues = [8, 7, 3, 4]
                    return (
                      <button 
                        key={label} 
                        type="button"
                        className="mood-chip-btn" 
                        style={{ '--chip-hue': moodColors[index] }}
                        onClick={() => handleQuickMood(label, moodValues[index], energyValues[index])}
                        aria-label={`Log ${label} mood`}
                      >
                        <span className="mood-chip-emoji" aria-hidden="true">{moodFaces[index]}</span>
                        <span className="mood-chip-label">{label}</span>
                      </button>
                    )
                  })}
                </section>

                <div className="metrics-dashboard-grid">
                  {/* Sleep Duration widget */}
                  <article className="metric-card sleep-metric-card">
                    <div className="metric-header">
                      <span className="metric-icon">🛌</span>
                      <h3>{t.sleepTitle}</h3>
                    </div>
                    <div className="metric-value">
                      <strong>{sleepHoursDisplay}</strong>
                    </div>
                    <div className="metric-chart">
                      {last6Entries.length === 0 ? (
                        <span className="chart-empty">{t.noEntriesYet}</span>
                      ) : (
                        <svg className="mini-bar-svg" viewBox="0 0 120 40">
                          {last6Entries.map((e, idx) => {
                            const val = e.sleep_hours || 0
                            const height = Math.max(2, Math.min(35, (val / 12) * 35))
                            const x = idx * 18 + 8
                            const y = 40 - height
                            return (
                              <rect
                                key={idx}
                                x={x}
                                y={y}
                                width="10"
                                height={height}
                                rx="3"
                                fill="var(--stress)"
                                opacity="0.85"
                              />
                            )
                          })}
                        </svg>
                      )}
                    </div>
                  </article>

                  {/* Stress Indicator widget */}
                  <article className="metric-card stress-metric-card">
                    <div className="metric-header">
                      <span className="metric-icon">📉</span>
                      <h3>{t.stressTitle}</h3>
                    </div>
                    <div className="metric-value">
                      <strong>{stressIndicatorLabel}</strong>
                    </div>
                    <div className="metric-chart">
                      {last6Entries.length === 0 ? (
                        <span className="chart-empty">{t.noEntriesYet}</span>
                      ) : (
                        <svg className="mini-sparkline-svg" viewBox="0 0 120 40">
                          {(() => {
                            const pointsStr = last6Entries.map((e, idx) => {
                              const x = idx * (120 / 5)
                              const stressVal = 11 - (e.mood_level || 6)
                              const y = 35 - (stressVal / 10) * 30
                              return `${x},${y}`
                            }).join(' ')
                            return (
                              <>
                                <polyline
                                  fill="none"
                                  stroke="var(--accent)"
                                  strokeWidth="2.5"
                                  points={pointsStr}
                                />
                                {last6Entries.map((e, idx) => {
                                  const x = idx * (120 / 5)
                                  const stressVal = 11 - (e.mood_level || 6)
                                  const y = 35 - (stressVal / 10) * 30
                                  return (
                                    <circle
                                      key={idx}
                                      cx={x}
                                      cy={y}
                                      r="3.5"
                                      fill="var(--accent)"
                                    />
                                  )
                                })}
                              </>
                            )
                          })()}
                        </svg>
                      )}
                    </div>
                  </article>

                  {/* Wellness Quiz Widget */}
                  <article className="metric-card quiz-metric-card">
                    <div className="metric-header">
                      <span className="metric-icon">⭐</span>
                      <h3>{t.quizTitle}</h3>
                    </div>
                    
                    <div className="quiz-body">
                      {quizStep < 4 ? (
                        <div className="quiz-step-content">
                          <p className="quiz-question-num">{t.questionNum.replace('{num}', quizStep + 1).replace('{total}', 4)}</p>
                          <p className="quiz-question-text">{quizQuestions[lang][quizStep]}</p>
                          <div className="quiz-buttons">
                            <button type="button" className="quiz-btn yes-btn" onClick={() => handleQuizAnswer(true)}>{t.yes}</button>
                            <button type="button" className="quiz-btn no-btn" onClick={() => handleQuizAnswer(false)}>{t.no}</button>
                          </div>
                        </div>
                      ) : (
                        <div className="quiz-result-content">
                          <p className="quiz-compl">{t.quizCompl}</p>
                          <p className="quiz-advice">{quizResultText}</p>
                          <button type="button" className="quiz-reset-btn" onClick={handleQuizReset}>{t.quizReset}</button>
                        </div>
                      )}
                    </div>
                  </article>
                </div>

                {/* Journal submit form */}
                <section className="panel journal-form-panel" aria-labelledby="form-title">
                  <h2 id="form-title">{t.writeJournal}</h2>
                  {status === 'error' && <p className="error-banner" role="alert">{error}</p>}
                  <form onSubmit={handleFullSubmit} className="journal-form">
                    <div className="form-group textarea-group">
                      <textarea
                        value={entryText}
                        onChange={(e) => setEntryText(e.target.value)}
                        placeholder={t.journalPlaceholder}
                        rows="4"
                        maxLength="4000"
                        required
                        aria-label={t.writeJournal}
                      />
                    </div>
                    
                    <div className="form-sliders-grid">
                      <div className="form-group slider-group">
                        <label htmlFor="mood-slider">
                          <span>{t.moodLevel}</span>
                          <strong className="slider-value">{moodLevel}</strong>
                        </label>
                        <input
                          type="range"
                          id="mood-slider"
                          min="1"
                          max="10"
                          value={moodLevel}
                          onChange={(e) => setMoodLevel(parseInt(e.target.value, 10))}
                        />
                      </div>

                      <div className="form-group slider-group">
                        <label htmlFor="energy-slider">
                          <span>{t.energyLevel}</span>
                          <strong className="slider-value">{energyLevel}</strong>
                        </label>
                        <input
                          type="range"
                          id="energy-slider"
                          min="1"
                          max="10"
                          value={energyLevel}
                          onChange={(e) => setEnergyLevel(parseInt(e.target.value, 10))}
                        />
                      </div>

                      <div className="form-group slider-group">
                        <label htmlFor="sleep-slider">
                          <span>{t.sleepDuration}</span>
                          <strong className="slider-value">{sleepHours}h</strong>
                        </label>
                        <input
                          type="range"
                          id="sleep-slider"
                          min="0"
                          max="24"
                          step="0.5"
                          value={sleepHours}
                          onChange={(e) => setSleepHours(parseFloat(e.target.value))}
                        />
                      </div>

                      <div className="form-group slider-group">
                        <label htmlFor="study-slider">
                          <span>{t.studyHours}</span>
                          <strong className="slider-value">{studyHours}h</strong>
                        </label>
                        <input
                          type="range"
                          id="study-slider"
                          min="0"
                          max="24"
                          step="0.5"
                          value={studyHours}
                          onChange={(e) => setStudyHours(parseFloat(e.target.value))}
                        />
                      </div>

                      <div className="form-group input-group">
                        <label htmlFor="countdown-input">{t.examCountdown}</label>
                        <input
                          type="number"
                          id="countdown-input"
                          min="0"
                          value={examCountdownDays}
                          onChange={(e) => setExamCountdownDays(Math.max(0, parseInt(e.target.value, 10) || 0))}
                        />
                      </div>
                    </div>

                    <div className="form-footer">
                      <button 
                        type="submit" 
                        className="primary-button submit-btn"
                        disabled={submitStatus === 'submitting'}
                      >
                        {submitStatus === 'submitting' ? t.submitting : t.submitEntry}
                      </button>
                    </div>
                  </form>
                </section>
              </div>
            )}

            {activeTab === 'calendar' && (
              <div className="tab-view calendar-view fade-in">
                <div className="calendar-panel-layout">
                  <section className="panel calendar-days-panel">
                    <header className="calendar-header-row">
                      <h2>June 2026</h2>
                      <span className="calendar-sub">{t.monthlyView}</span>
                    </header>

                    <div className="calendar-grid-wrapper">
                      <div className="calendar-week-days">
                        <span>Sun</span>
                        <span>Mon</span>
                        <span>Tue</span>
                        <span>Wed</span>
                        <span>Thu</span>
                        <span>Fri</span>
                        <span>Sat</span>
                      </div>
                      
                      <div className="calendar-days-grid" role="grid" aria-label="Mood calendar for June 2026">
                        {calendarDays.map((day, idx) => {
                          const entryForDay = entries.find(e => e.created_at && e.created_at.startsWith(day.date))
                          const dayMoodState = entryForDay ? getMoodState(entryForDay.mood_level || 5, lang) : null
                          
                          return (
                            <button
                              key={`${day.date}-${idx}`}
                              type="button"
                              className={`calendar-cell ${day.isCurrentMonth ? 'current-month' : 'other-month'} ${entryForDay ? 'has-entry' : ''}`}
                              onClick={() => entryForDay && handleEntryClick(entryForDay)}
                              disabled={!entryForDay}
                              style={dayMoodState ? { '--cell-color': dayMoodState.color } : {}}
                              aria-label={entryForDay ? `${day.date}, logged mood: ${dayMoodState.label}` : `${day.date}, no entry`}
                            >
                              <span className="day-number">{day.dayNum}</span>
                              {dayMoodState && (
                                <span className="day-emoji" aria-hidden="true">{dayMoodState.face}</span>
                              )}
                            </button>
                          )
                        })}
                      </div>
                    </div>
                  </section>

                  {/* Monthly Mood Summary */}
                  <div className="calendar-sidebar">
                    <section className="panel dominant-mood-panel" style={dominantMood ? { background: `color-mix(in srgb, ${dominantMood.color} 18%, var(--surface))` } : {}}>
                      <p className="section-kicker">{t.monthlySummary}</p>
                      <h2>{dominantMood ? dominantMood.label : (lang === 'hi' ? 'कोई डेटा नहीं' : 'No Data')}</h2>
                      
                      <div className="dominant-mood-body">
                        <span className="dominant-mood-emoji" aria-hidden="true">{dominantMood ? dominantMood.face : '😶'}</span>
                        <p className="dominant-mood-desc">{dominantMoodDesc}</p>
                      </div>
                    </section>

                    {/* Supplementary metrics */}
                    <section className="panel calendar-metrics-panel">
                      <h3>Tracking Metrics</h3>
                      <div className="supplementary-metrics-list">
                        <div className="supp-metric-item">
                          <span className="supp-icon">🏃</span>
                          <div className="supp-info">
                            <span>{t.activity}</span>
                            <strong>8,432 {t.steps}</strong>
                          </div>
                        </div>
                        <div className="supp-metric-item">
                          <span className="supp-icon">💬</span>
                          <div className="supp-info">
                            <span>{t.therapy}</span>
                            <strong>2 {t.sessions}</strong>
                          </div>
                        </div>
                        <div className="supp-metric-item">
                          <span className="supp-icon">📚</span>
                          <div className="supp-info">
                            <span>{t.discipline}</span>
                            <strong>88% {t.focusScore}</strong>
                          </div>
                        </div>
                      </div>
                    </section>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'trends' && (
              <div className="tab-view trends-view fade-in">
                <section className="panel insight-panel" aria-labelledby="insight-title">
                  <p className="section-kicker">{t.todaysSnapshot}</p>
                  <h2 id="insight-title">{t.quickInsight}</h2>
                  <div className="insight-face" aria-hidden="true">
                    {latestMoodState ? latestMoodState.face : '🙂'}
                  </div>
                  <p className="insight-copy">
                    {latestPoint
                      ? t.steadyMoodMsg
                          .replace('{moodLabel}', lang === 'en' ? (latestMoodState?.label?.toLowerCase() || 'steady') : (latestMoodState?.label || 'steady'))
                          .replace('{stress}', formatValue(latestPoint.stress_average))
                          .replace('{mood}', formatValue(latestPoint.mood_average))
                      : t.noInsightMsg}
                  </p>
                  <div className="mini-metrics">
                    <div>
                      <span>{t.moodLabel}</span>
                      <strong>{latestPoint ? formatValue(latestPoint.mood_average) : '0.0'}</strong>
                    </div>
                    <div>
                      <span>{t.energy}</span>
                      <strong>{latestPoint ? formatValue(latestPoint.energy_average) : '0.0'}</strong>
                    </div>
                  </div>
                </section>

                <TrendChart points={points} lang={lang} />

                <section className="panel entries-panel" aria-labelledby="entries-title">
                  <div className="section-heading">
                    <p className="section-kicker">{t.recentEntries}</p>
                    <h2 id="entries-title">{t.savedNotes}</h2>
                  </div>
                  <div className="entry-list">
                    {entries.length === 0 ? (
                      <p className="empty-state">{t.noNotesYet}</p>
                    ) : (
                      entries.slice(-6).reverse().map((entry) => {
                        const entryMood = getMoodState(entry.mood_level || entry.mood_average || 5, lang)
                        return (
                          <article key={entry.id} className="entry-card" onClick={() => handleEntryClick(entry)} style={{ cursor: 'pointer' }}>
                            <div className="entry-meta">
                              <span className="entry-date">{formatShortDate(entry.created_at)}</span>
                              <span className="entry-face" aria-hidden="true">{entryMood.face}</span>
                            </div>
                            <p>{entry.entry_text}</p>
                            <div className="entry-foot">
                              <span>{t.moodLabel} {entry.mood_level}</span>
                              <span>{t.stressLabel} {formatValue(11 - entry.mood_level)}</span>
                            </div>
                          </article>
                        )
                      })
                    )}
                  </div>
                </section>
              </div>
            )}

            {activeTab === 'coping' && (
              <div className="tab-view coping-view fade-in">
                {copingData ? (
                  <section className={`panel companion-panel ${copingData.is_crisis ? 'crisis-mode' : ''}`} aria-labelledby="companion-title">
                    <div className="section-heading">
                      <p className="section-kicker">{copingData.is_crisis ? t.crisisKicker : t.comfortKicker}</p>
                      <h2 id="companion-title">{copingData.is_crisis ? t.crisisTitle : t.comfortTitle}</h2>
                    </div>
                    
                    {copingData.is_crisis ? (
                      <div className="crisis-content">
                        <div className="crisis-alert-icon" aria-hidden="true">⚠️</div>
                        <p className="crisis-warning"><strong>{t.crisisWarning}</strong> {copingData.guidance.motivational_prompt}</p>
                        
                        <div className="helpline-box">
                          <strong>{t.aasra}</strong>
                          <a href="tel:+919820466726" className="helpline-link">+91 98204 66726</a>
                        </div>
                        
                        <div className="helpline-box">
                          <strong>{t.vandrevala}</strong>
                          <a href="tel:9999666555" className="helpline-link">9999 666 555</a>
                        </div>

                        <p className="crisis-subtext">{copingData.guidance.mindfulness_activity}</p>
                      </div>
                    ) : (
                      <div className="companion-content">
                        <div className="companion-tabs" role="tablist">
                          <button 
                            type="button" 
                            role="tab" 
                            aria-selected={copingTab === 'motivation'}
                            className={`tab-btn ${copingTab === 'motivation' ? 'active' : ''}`}
                            onClick={() => setCopingTab('motivation')}
                          >
                            {t.comfort}
                          </button>
                          <button 
                            type="button" 
                            role="tab" 
                            aria-selected={copingTab === 'breathing'}
                            className={`tab-btn ${copingTab === 'breathing' ? 'active' : ''}`}
                            onClick={() => setCopingTab('breathing')}
                          >
                            {t.breathe}
                          </button>
                          <button 
                            type="button" 
                            role="tab" 
                            aria-selected={copingTab === 'mindfulness'}
                            className={`tab-btn ${copingTab === 'mindfulness' ? 'active' : ''}`}
                            onClick={() => setCopingTab('mindfulness')}
                          >
                            {t.grounding}
                          </button>
                        </div>

                        <div className="tab-body" role="tabpanel">
                          {copingTab === 'motivation' && (
                            <div className="tab-pane fade-in">
                              <span className="quote-mark">“</span>
                              <p className="companion-motivation-text">{copingData.guidance.motivational_prompt}</p>
                            </div>
                          )}
                          {copingTab === 'breathing' && (
                            <div className="tab-pane fade-in">
                              <h4>{t.breathingPractice}</h4>
                              <p className="companion-exercise-text">{copingData.guidance.breathing_exercise}</p>
                            </div>
                          )}
                          {copingTab === 'mindfulness' && (
                            <div className="tab-pane fade-in">
                              <h4>{t.mindfulGrounding}</h4>
                              <p className="companion-mindfulness-text">{copingData.guidance.mindfulness_activity}</p>
                            </div>
                          )}
                        </div>
                      </div>
                    )}
                  </section>
                ) : (
                  <div className="panel empty-coping-panel">
                    <p>{t.noInsightMsg}</p>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Bottom glassmorphic navigation bar */}
          <nav className="bottom-nav-bar" aria-label="Main Navigation">
            <button
              type="button"
              className={`nav-item ${activeTab === 'home' ? 'active' : ''}`}
              onClick={() => setActiveTab('home')}
            >
              <span className="nav-icon" aria-hidden="true">🏠</span>
              <span className="nav-label">{t.tabHome}</span>
            </button>
            <button
              type="button"
              className={`nav-item ${activeTab === 'calendar' ? 'active' : ''}`}
              onClick={() => setActiveTab('calendar')}
            >
              <span className="nav-icon" aria-hidden="true">📅</span>
              <span className="nav-label">{t.tabCalendar}</span>
            </button>
            <button
              type="button"
              className={`nav-item ${activeTab === 'trends' ? 'active' : ''}`}
              onClick={() => setActiveTab('trends')}
            >
              <span className="nav-icon" aria-hidden="true">📊</span>
              <span className="nav-label">{t.tabTrends}</span>
            </button>
            <button
              type="button"
              className={`nav-item ${activeTab === 'coping' ? 'active' : ''}`}
              onClick={() => setActiveTab('coping')}
            >
              <span className="nav-icon" aria-hidden="true">❤️</span>
              <span className="nav-label">{t.tabCoping}</span>
            </button>
          </nav>
        </>
      )}

      {selectedEntry && (
        <div className="modal-overlay" onClick={() => setSelectedEntry(null)} role="dialog" aria-modal="true" aria-labelledby="modal-title">
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <header className="modal-header">
              <div>
                <p className="section-kicker">{t.stressTriggerAnalysis}</p>
                <h2 id="modal-title">{t.journalInsights}</h2>
              </div>
              <button type="button" className="close-button" onClick={() => setSelectedEntry(null)} aria-label={t.closeModal}>
                &times;
              </button>
            </header>
            
            <div className="modal-body">
              <div className="entry-context-card">
                <div className="entry-context-meta">
                  <span className="entry-date">{formatShortDate(selectedEntry.created_at)}</span>
                  <span className="entry-face" aria-hidden="true">
                    {getMoodState(selectedEntry.mood_level || 5, lang).face} {t.moodLabel} {selectedEntry.mood_level}/10
                  </span>
                </div>
                <p className="entry-context-text">"${selectedEntry.entry_text}"</p>
              </div>

              {analysisLoading && (
                <div className="modal-loader">
                  <div className="spinner"></div>
                  <p>{t.analyzingProgress}</p>
                </div>
              )}

              {analysisError && (
                <div className="error-banner" role="alert">
                  {analysisError}
                </div>
              )}

              {analysis && (
                <div className="analysis-results">
                  <div className="analysis-metrics-row">
                    <div className="metric-box stress-meter-box">
                      <span>{t.stressScore}</span>
                      <strong className="stress-value">{formatValue(analysis.stress_score)}</strong>
                      <div className="stress-bar-wrapper">
                        <div 
                          className="stress-bar-fill" 
                          style={{ 
                            width: `${analysis.stress_score * 10}%`,
                            backgroundColor: analysis.stress_score > 7 ? 'var(--stress-deep)' : analysis.stress_score > 4 ? 'var(--stress)' : '#92e38f'
                          }}
                        ></div>
                      </div>
                      <small className="scale-label">{t.scaleLabel}</small>
                    </div>

                    <div className="metric-box study-sleep-box">
                      <div className="sub-metric">
                        <span>{t.sleep}</span>
                        <strong>{selectedEntry.sleep_hours}h</strong>
                      </div>
                      <div className="sub-metric">
                        <span>{t.study}</span>
                        <strong>{selectedEntry.study_hours}h</strong>
                      </div>
                      <div className="sub-metric">
                        <span>{t.countdown}</span>
                        <strong>{selectedEntry.exam_countdown_days}d</strong>
                      </div>
                    </div>
                  </div>

                  <div className="analysis-section">
                    <h3>{t.extractedTriggers}</h3>
                    {analysis.triggers && analysis.triggers.length > 0 ? (
                      <div className="trigger-chips-list">
                        {analysis.triggers.map((trigger, idx) => (
                          <span key={idx} className="trigger-chip">
                            ⚠️ {trigger}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <p className="no-triggers-text">{t.noTriggers}</p>
                    )}
                  </div>

                  <div className="analysis-section summary-box">
                    <h3>{t.wellnessSummary}</h3>
                    <p className="analysis-summary-text">{analysis.summary}</p>
                  </div>

                  <div className="analysis-section explanation-box">
                    <h3>{t.empatheticGuidance}</h3>
                    <p className="analysis-explanation-text">{analysis.empathetic_explanation}</p>
                  </div>
                </div>
              )}
            </div>

            <footer className="modal-footer">
              <button type="button" className="primary-button" onClick={() => setSelectedEntry(null)}>
                {t.gotIt}
              </button>
            </footer>
          </div>
        </div>
      )}

      {toastMessage && (
        <div className={`toast-alert toast-${submitStatus} fade-in`} role="alert">
          <span className="toast-icon">{submitStatus === 'success' ? '✓' : '⚠️'}</span>
          <span className="toast-text">{toastMessage}</span>
        </div>
      )}

    </main>
  )
}

export default App;
