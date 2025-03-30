'use client';

import React, { ReactNode } from 'react';

interface SelectProps {
  children: ReactNode;
  value?: string;
  onValueChange?: (value: string) => void;
  name?: string;
}

interface SelectItemProps {
  value: string;
  children: React.ReactNode;
  ["data-value"]?: string;
  role?: string;
  ["aria-selected"]?: boolean;
}

interface SelectContentProps {
  children: ReactNode;
}

// Define a more specific type for props we might encounter
type ChildProps = {
  children?: ReactNode;
  value?: string;
  [key: string]: unknown; // Use unknown instead of any for better type safety
};

const SelectTrigger = ({ children, ...props }: { children: ReactNode } & React.HTMLAttributes<HTMLDivElement>) => (
  <div data-testid="mock-select-trigger" {...props}>{children}</div>
);

const SelectValue = ({ placeholder }: { placeholder?: string }) => (
  <span data-testid="mock-select-value">{placeholder}</span>
);

const SelectContent = ({ children }: SelectContentProps) => (
  <>{children}</>
);
SelectContent.displayName = 'SelectContent';

// Add a minimal SelectGroup mock
const SelectGroup = ({ children }: { children: ReactNode }) => <>{children}</>;
SelectGroup.displayName = 'SelectGroup';

// Add minimal mocks for scroll buttons and viewport
const SelectScrollUpButton = ({ children }: { children?: ReactNode }) => <>{children}</>;
SelectScrollUpButton.displayName = 'SelectScrollUpButton';

const SelectScrollDownButton = ({ children }: { children?: ReactNode }) => <>{children}</>;
SelectScrollDownButton.displayName = 'SelectScrollDownButton';

const SelectViewport = ({ children }: { children?: ReactNode }) => <>{children}</>;
SelectViewport.displayName = 'SelectViewport';

// Add minimal SelectLabel mock
const SelectLabel = ({ children }: { children: ReactNode }) => <>{children}</>;
SelectLabel.displayName = 'SelectLabel';

// Add minimal SelectSeparator mock
const SelectSeparator = () => <hr data-testid="mock-separator" />;
SelectSeparator.displayName = 'SelectSeparator';

const Select = React.forwardRef<HTMLSelectElement, SelectProps>(({ value, onValueChange, children, name, ...rest }, ref) => {
  const options: { value: string; label: ReactNode }[] = [];

  const findItems = (nodes: ReactNode) => {
    React.Children.forEach(nodes, (child) => {
      if (React.isValidElement(child)) {
        const childType = child.type as React.ComponentType & { displayName?: string };
        const props = child.props as ChildProps;

        const isSelectItem = childType.displayName === 'SelectItem' || 'value' in props;
        const isSelectContent = childType.displayName === 'SelectContent';

        if (isSelectItem && props.value !== undefined) {
          options.push({
            value: String(props.value || props["data-value"] || ''), // Ensure value is string
            label: props.children,
          });
        } else if (isSelectContent || (typeof childType !== 'string' && props.children)) {
          findItems(props.children);
        }
      }
    });
  };

  findItems(children);

  return (
    <select
      ref={ref}
      data-testid="native-select"
      value={value ?? ''}
      onChange={(e) => onValueChange?.(e.target.value)}
      name={name}
      id="status-select"
      aria-label="Filter by Status"
      {...rest}
    >
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  );
});
Select.displayName = 'Select';

const SelectItem = ({ children }: SelectItemProps) => <>{children}</>;
SelectItem.displayName = 'SelectItem';

export {
  SelectContent as Content,
  SelectGroup as Group,
  SelectItem as Item,
  SelectLabel as Label,
  Select as Root,
  SelectScrollDownButton as ScrollDownButton,
  SelectScrollUpButton as ScrollUpButton,
  SelectSeparator as Separator,
  SelectTrigger as Trigger,
  SelectValue as Value,
  SelectViewport as Viewport
};
