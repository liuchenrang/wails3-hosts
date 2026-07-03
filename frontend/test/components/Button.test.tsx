import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from '../../src/components/ui/Button';

describe('Button 组件', () => {
  it('应该正确渲染按钮', () => {
    render(<Button>点击我</Button>);
    expect(screen.getByRole('button')).toHaveTextContent('点击我');
  });

  it('应该支持不同的变体样式', () => {
    const { rerender } = render(<Button variant="default">默认</Button>);
    expect(screen.getByRole('button')).toBeInTheDocument();

    rerender(<Button variant="outline">轮廓</Button>);
    expect(screen.getByRole('button')).toBeInTheDocument();

    rerender(<Button variant="ghost">幽灵</Button>);
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('应该支持禁用状态', () => {
    render(<Button disabled>禁用</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('应该支持 onClick 事件', async () => {
    const handleClick = vi.fn();
    const { userEvent } = await import('@testing-library/user-event');
    const user = userEvent.setup();

    render(<Button onClick={handleClick}>点击</Button>);
    await user.click(screen.getByRole('button'));

    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});