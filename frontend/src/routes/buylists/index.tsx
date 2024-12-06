import { createFileRoute, useNavigate } from "@tanstack/react-router";

import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ListBuylists } from "@wails/go/main/App.js";

import { NewBuylistDialog } from "@/components/buylists/newBuylist";

export const Route = createFileRoute("/buylists/")({
  component: Index,
  loader: () => ListBuylists(),
  errorComponent: ({ error, reset }) => {
    return (
      <div>
        <pre>
          ({typeof error}) {JSON.stringify(error, null, 2)}
        </pre>
        <button
          onClick={() => {
            reset();
          }}
        >
          retry
        </button>
      </div>
    );
  },
});

function Index() {
  const buylists = Route.useLoaderData();
  const navigate = useNavigate();
  return (
    <div className="container mx-auto pt-8">
      <div className="flex flex-row pb-4">
        <h2 className="scroll-m-20 text-xl font-semibold tracking-tight pb-1">
          Buylists
        </h2>
        <div className="flex gap-2 ml-auto">
          <NewBuylistDialog />
        </div>
      </div>
      <Table>
        <TableCaption>Select a buylist to get started</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>Buylist Name</TableHead>
            <TableHead>Created</TableHead>
            <TableHead className="text-right">Total Card Count</TableHead>
            <TableHead className="text-right">Unique Card Count</TableHead>
            <TableHead className="text-right">Total Selection Cost</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {buylists.map((bl) => {
            return (
              <TableRow
                key={bl.id}
                onClick={() => {
                  void navigate({
                    to: "/buylists/$buylistId",
                    params: { buylistId: String(bl.id) },
                  });
                }}
                className="cursor-pointer"
              >
                <TableCell>{bl.name}</TableCell>
                <TableCell>{new Date(bl.createdAt).toISOString()}</TableCell>
                <TableCell className="text-right">{bl.totalCards}</TableCell>
                <TableCell className="text-right">{bl.cards.length}</TableCell>
                <TableCell className="text-right">
                  ${String(bl.totalPrice / 100)}
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
