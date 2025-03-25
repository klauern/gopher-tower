# Shadcn/UI Migration Plan

## Overview

This document outlines the step-by-step plan for migrating our frontend components to use shadcn/ui. The migration will be done incrementally to ensure stability and maintain functionality throughout the process.

## Prerequisites ✅

- [x] Initialize shadcn/ui with `shadcn-ui init`
- [x] Setup base configuration in `components.json`
- [x] Configure tailwind.config.ts for shadcn/ui

## 1. Component Audit and Planning

### Current Component Structure

- `/components`
  - `EventStream.tsx` - Complex component for handling event streams
  - `Navbar.tsx` - Navigation component
  - `/ui` - Base UI components
  - `/jobs` - Job-related components

### Shadcn/UI Components to Add

#### Priority 1 - Core Components ✅

- [x] Button (`@/components/ui/button`)
- [x] Input (`@/components/ui/input`)
- [x] Form (`@/components/ui/form`)
- [x] Card (`@/components/ui/card`)
- [x] Dialog (`@/components/ui/dialog`)
- [x] Toast (`@/components/ui/toast`)
- [x] Navigation Menu (`@/components/ui/navigation-menu`)

#### Priority 2 - Enhanced Components

- [ ] Table (`@/components/ui/table`)
- [ ] Tabs (`@/components/ui/tabs`)
- [ ] Sheet (`@/components/ui/sheet`)
- [ ] Select (`@/components/ui/select`)
- [ ] Command (`@/components/ui/command`)

## 2. Migration Steps

### Phase 1: Core Infrastructure

1. [x] Setup shadcn/ui theme configuration
   - [x] Review and customize `/lib/utils.ts`
   - [x] Configure dark mode support
   - [x] Define color palette variables (using shadcn/ui default palette)

2. [x] Add base components

   ```bash
   npx shadcn@latest add button input form card dialog toast navigation-menu dropdown-menu
   ```

### Phase 2: Component Migration

#### Navbar Migration

1. [ ] Replace current Navbar with shadcn/ui navigation menu
   - Migrate existing links and structure
   - Implement responsive design
   - Add dropdown menus if needed

#### EventStream Component

1. [ ] Break down into smaller components
2. [ ] Integrate shadcn/ui components:
   - Use Card for event containers
   - Add Toast for notifications
   - Implement Dialog for detailed views

#### Jobs Components

1. [ ] Audit existing job components
2. [ ] Plan component replacements
3. [ ] Implement new components using shadcn/ui

### Phase 3: Enhanced Features

1. [x] Add dark mode toggle
2. [ ] Implement toast notifications system
3. [ ] Add loading states and animations
4. [ ] Enhance form validations

## 3. Testing Strategy

### Unit Tests

- [ ] Update existing component tests
- [ ] Add new tests for shadcn/ui components
- [ ] Verify dark mode functionality
- [ ] Test responsive behavior

### Integration Tests

- [ ] Test component interactions
- [ ] Verify form submissions
- [ ] Test navigation flows
- [ ] Validate event handling

## 4. Documentation

1. [ ] Update component documentation
2. [ ] Add usage examples
3. [ ] Document theme customization
4. [ ] Update TESTING.md with new testing patterns

## 5. Performance Considerations

- [ ] Analyze bundle size impact
- [ ] Implement code splitting
- [ ] Optimize component loading
- [ ] Monitor performance metrics

## Progress Tracking

Use this section to track progress as we implement changes:

```markdown
### Current Status
- Phase: Phase 1 - Core Infrastructure
- Next Action: Begin Phase 2 - Component Migration (starting with Navbar)
- Completed Items: 10/25
```

## Notes and Decisions

- All new components should be client components ('use client')
- Maintain TypeScript strict mode
- Follow existing testing patterns
- Keep accessibility features intact
