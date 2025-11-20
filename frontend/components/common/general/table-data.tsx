import * as React from "react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { typeConfig, statusConfig } from "@/components/common/general/table-filter"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table"
import { Spinner } from "@/components/ui/spinner"
import { ErrorInline } from "@/components/layout/error"
import { EmptyStateWithBorder } from "@/components/layout/empty"
import { LoadingStateWithBorder } from "@/components/layout/loading"
import { ListRestart, Layers, LucideIcon } from "lucide-react"
import { formatDateTime } from "@/lib/utils"
import type { Order } from "@/lib/services"


/**
 * 交易数据表格组件
 * 显示交易记录的表格，可复用于不同页面
 */
export function TransactionDataTable({ transactions }: { transactions: Order[] }) {
  return (
    <div className="border border-dashed shadow-none rounded-lg overflow-hidden">
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="border-b border-dashed">
              <TableHead className="whitespace-nowrap w-[120px]">名称</TableHead>
              <TableHead className="whitespace-nowrap text-center min-w-[160px]">订单号</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[120px]">商户订单号</TableHead>
              <TableHead className="whitespace-nowrap text-right min-w-[80px]">金额</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[80px]">类型</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[80px]">状态</TableHead>
              <TableHead className="whitespace-nowrap text-center min-w-[100px]">交易双方</TableHead>
              <TableHead className="whitespace-nowrap text-center min-w-[50px]">商户</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[120px]">创建时间</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[120px]">过期时间</TableHead>
              <TableHead className="whitespace-nowrap text-center w-[120px]">交易时间</TableHead>
              <TableHead className="sticky right-0 whitespace-nowrap text-center bg-background shadow-[-4px_0_8px_-2px_rgba(0,0,0,0.1)] w-[150px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="animate-in fade-in duration-200">
            {transactions.map((order) => (
              <TransactionTableRow
                key={order.order_no}
                order={order}
              />
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}

/**
 * 交易表格行组件
 * 展示交易记录的表格行
 */
function TransactionTableRow({ order }: { order: Order }) {
  /* 格式化金额 */
  const getAmountDisplay = (amount: string) => (
    <span className="text-xs font-semibold">
      {parseFloat(amount).toFixed(2)}
    </span>
  )

  return (
    <TableRow className="h-8 border-b border-dashed">
      <TableCell className="text-xs font-medium whitespace-nowrap py-1">
        {order.order_name}
      </TableCell>
      <TableCell className="font-mono text-xs text-center py-1">
        {order.order_no}
      </TableCell>
      <TableCell className="font-mono text-xs text-center py-1">
        {order.merchant_order_no || '-'}
      </TableCell>
      <TableCell className="whitespace-nowrap text-right py-1">
        {getAmountDisplay(order.amount)}
      </TableCell>
      <TableCell className="whitespace-nowrap text-center py-1">
        <Badge
          variant="secondary"
          className={`text-[11px] px-1 ${typeConfig[order.type].color}`}
        >
          {typeConfig[order.type].label}
        </Badge>
      </TableCell>
      <TableCell className="whitespace-nowrap text-center py-1">
        <Badge
          variant="secondary"
          className={`text-[11px] px-1 ${statusConfig[order.status].color}`}
        >
          {statusConfig[order.status].label}
        </Badge>
      </TableCell>
      <TableCell className="whitespace-nowrap text-center py-1">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="flex items-center cursor-pointer gap-1 justify-center">
                <Avatar className="h-4 w-4">
                  <AvatarImage src={undefined} />
                  <AvatarFallback className="text-[9px] bg-primary text-primary-foreground">
                    {order.payer_username.substring(0, 1).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div className="text-xs font-bold">⭢</div>
                <Avatar className="h-4 w-4">
                  <AvatarImage src={undefined} />
                  <AvatarFallback className="text-[9px] bg-primary text-primary-foreground">
                    {order.payee_username.substring(0, 1).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
              </div>
            </TooltipTrigger>
            <TooltipContent side="top" className="p-3">
              <div className="space-y-2">
                <div>
                  <p className="text-xs font-semibold">付款方</p>
                  <p className="text-xs">ID: {order.payer_user_id}</p>
                  <p className="text-xs">账户: {order.payer_username}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold">收款方</p>
                  <p className="text-xs">ID: {order.payee_user_id}</p>
                  <p className="text-xs">账户: {order.payee_username}</p>
                </div>
              </div>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </TableCell>
      <TableCell className="text-xs text-center py-1">
        {order.app_name || '-'}
      </TableCell>
      <TableCell className="text-xs text-center py-1">
        {formatDateTime(order.created_at)}
      </TableCell>
      <TableCell className="text-xs text-center py-1">
        {formatDateTime(order.expires_at)}
      </TableCell>
      <TableCell className="text-xs text-center py-1">
        {(order.status === 'success' || order.status === 'refund') ? formatDateTime(order.trade_time) : '-'}
      </TableCell>
      <TableCell className="sticky right-0 whitespace-nowrap text-center bg-background shadow-[-4px_0_8px_-2px_rgba(0,0,0,0.1)] py-1">
        <Button variant="ghost" size="sm" className="h-6 px-2 text-xs">
          详情
        </Button>
      </TableCell>
    </TableRow>
  )
}

interface TransactionTableListProps {
  loading: boolean
  error: Error | null
  transactions: Order[]
  total: number
  currentPage: number
  totalPages: number
  onRetry: () => void
  onLoadMore: () => void
  emptyIcon?: LucideIcon
  emptyDescription?: string
}

/**
 * 交易列表容器组件
 * 统一处理加载、错误、空状态和分页加载
 */
export function TransactionTableList({
  loading,
  error,
  transactions,
  total,
  currentPage,
  totalPages,
  onRetry,
  onLoadMore,
  emptyIcon = Layers,
  emptyDescription = "未发现活动"
}: TransactionTableListProps) {
  if (loading && transactions.length === 0) {
    return (
      <LoadingStateWithBorder
        icon={ListRestart}
        description="数据加载中"
      />
    )
  }

  if (error) {
    return (
      <div className="p-8 border-2 border-dashed border-border rounded-lg">
        <ErrorInline
          error={error}
          onRetry={onRetry}
          className="justify-center"
        />
      </div>
    )
  }

  if (!transactions || transactions.length === 0) {
    return (
      <EmptyStateWithBorder
        icon={emptyIcon}
        description={emptyDescription}
      />
    )
  }

  return (
    <div className="flex flex-col gap-2">
      <TransactionDataTable transactions={transactions} />

      {currentPage < totalPages && (
        <Button
          variant="outline"
          onClick={onLoadMore}
          disabled={loading}
          className="w-full text-xs border-dashed shadow-none"
        >
          {loading ? (<><Spinner className="size-4" />正在加载</>) : (`加载更多 (${transactions.length}/${total})`)}
        </Button>
      )}

      {currentPage >= totalPages && total > 0 && (
        <div className="pt-2 text-center text-xs text-muted-foreground">
          已加载全部 {total} 条记录
        </div>
      )}
    </div>
  )
}
