import { Section } from "@/components/layout/section";
import { Price } from "@/components/price/priceDisplay";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import {
  useSelectionRemoveCallback,
  useSelectionStatusCallback,
} from "@/hooks/selectionsStatus";
import { createFileRoute } from "@tanstack/react-router";
import { GetBuylistSummary } from "@wails/go/main/App";
import { buylist, models } from "@wails/go/models";
import { BrowserOpenURL } from "@wails/runtime/runtime";
import { Check, ExternalLink, Flag, Trash2, TriangleAlert } from "lucide-react";

export const Route = createFileRoute("/buylists/$buylistId/checkout")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const summary = await GetBuylistSummary(Number(params.buylistId));
    return summary;
  },
});

function RouteComponent() {
  const { summaryByVendor, vendors } = Route.useLoaderData();
  return (
    <div className="m-4">
      <Alert className="mb-4">
        <TriangleAlert />
        <div className="ml-2">
          <AlertTitle>Heads up!</AlertTitle>
          <AlertDescription>
            Pricing and quantity data is a best-effort approximation and may
            change during checkout. Please double check when checking out with
            your vendor.
          </AlertDescription>
        </div>
      </Alert>
      <Section header="Buylist Summary">
        <Tabs defaultValue={vendors[0]}>
          <TabsList>
            {vendors.map((vendor) => (
              <TabsTrigger key={vendor} value={vendor}>
                {vendor}
              </TabsTrigger>
            ))}
          </TabsList>
          {vendors.map((vendor) => {
            return (
              <TabsContent value={vendor} key={vendor}>
                <SummaryTable summary={summaryByVendor[vendor]} />
              </TabsContent>
            );
          })}
        </Tabs>
      </Section>
    </div>
  );
}

interface SummaryTableProps {
  summary: models.VendorSummary;
}
const SummaryTable = ({ summary }: SummaryTableProps) => {
  const exportText = summary.cardList
    .map((cl) => `${cl.qty.toString()} ${cl.name}`)
    .join("\n");
  const updateSelection = useSelectionStatusCallback();
  const removeSelection = useSelectionRemoveCallback();
  return (
    <div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Item Name</TableHead>
            <TableHead>Condition</TableHead>

            <TableHead className="text-right">Cost/Item</TableHead>
            <TableHead className="text-right">Qty Available</TableHead>

            <TableHead className="text-right">Quantity</TableHead>
            <TableHead className="text-right">Total Cost</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {summary.selections.map((sel) => (
            <TableRow key={sel.id}>
              <TableCell>
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    BrowserOpenURL(sel.offering.productUri);
                  }}
                  className="flex flex-row space-x-1"
                >
                  <span>{sel.offering.name}</span>
                  <span>
                    <ExternalLink size={18} />
                  </span>
                </a>
              </TableCell>
              <TableCell>{sel.offering.condition}</TableCell>

              <TableCell className="text-right">
                <Price price={sel.offering.price} />
              </TableCell>
              <TableCell className="text-right">
                {sel.offering.quantity} available
              </TableCell>
              <TableCell className="text-right">{sel.quantity}</TableCell>
              <TableCell className="text-right">
                <Price price={sel.offering.price * sel.quantity} />
              </TableCell>
              <TableCell>
                <Button
                  size="icon"
                  variant={sel.isPurchased ? "default" : "ghost"}
                  onClick={() => {
                    updateSelection(
                      new buylist.UpdateSelectionRequest({
                        selectionId: sel.id,
                        isPurchased: !sel.isPurchased,
                      }),
                    );
                  }}
                >
                  <Check />
                </Button>
                <Button
                  size="icon"
                  variant={sel.isFlagged ? "destructive" : "ghost"}
                  onClick={() => {
                    updateSelection(
                      new buylist.UpdateSelectionRequest({
                        selectionId: sel.id,
                        isFlagged: !sel.isFlagged,
                      }),
                    );
                  }}
                >
                  <Flag />
                </Button>
                <Button
                  size="icon"
                  variant="ghost"
                  onClick={() => {
                    removeSelection(sel);
                  }}
                >
                  <Trash2 />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
        <TableFooter>
          <TableRow>
            <TableCell colSpan={4}>Total</TableCell>
            <TableCell className="text-right">{summary.totalQty}</TableCell>
            <TableCell className="text-right">
              <Price price={summary.totalCost} />
            </TableCell>
          </TableRow>
        </TableFooter>
      </Table>

      <div className="pt-4 grid grid-cols-2">
        <div className="p-4">
          <Textarea
            className="resize-none"
            rows={20}
            value={exportText}
            contentEditable={false}
          />
        </div>
        <div className="p-4">
          <Textarea
            className="resize-none"
            rows={20}
            value={"TODO: cart automation stuff here"}
            contentEditable={false}
          />
        </div>
      </div>
    </div>
  );
};
