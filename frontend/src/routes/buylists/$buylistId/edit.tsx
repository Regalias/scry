import { createFileRoute, useRouter } from "@tanstack/react-router";
import {
  AddCardsToBuylist,
  DeleteBuylist,
  GetBuylist,
  UpdateBuylistName,
} from "@wails/go/main/App";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { ButtonWithConfirm } from "@/components/buttons/buttonWithConfirm";
import { EditStringForm } from "@/components/forms/editString";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "@/hooks/use-toast";
import { buylist } from "@wails/go/models";

export const Route = createFileRoute("/buylists/$buylistId/edit")({
  component: RouteComponent,
  loader: async (ctx) => {
    const id = Number(ctx.params.buylistId);
    return await GetBuylist(id);
  },
});

// Form schemas
const formSchema = z.object({
  cardlist: z.string().min(1, {
    message: "Must not be blank",
  }),
});
type schemaType = z.infer<typeof formSchema>;

// Regex for decklist parser
const re = new RegExp("^([\\d]+(?:x)?) ([\\S ]+)$");

// Component
function RouteComponent() {
  const buylist = Route.useLoaderData();
  const navigate = Route.useNavigate();
  const router = useRouter();

  const form = useForm<schemaType>({
    resolver: zodResolver(formSchema),
  });

  function onSubmit(data: schemaType) {
    const parsedCards: buylist.AddCardsRequest[] = [];
    for (const line of data.cardlist.split("\n")) {
      if (!line) continue;

      const result = re.exec(line.trim());
      if (result?.length != 3) {
        form.setError("cardlist", {
          message: `Line had an invalid format: '${line.trim()}'`,
        });
        return;
      }
      parsedCards.push({ name: result[2], quantity: Number(result[1]) });
    }

    AddCardsToBuylist(buylist.id, parsedCards)
      .then(() => {
        toast({
          description: `Added ${String(parsedCards.length)} unique card(s) to the buylist`,
        });
        void router.invalidate();
      })
      .catch((err: unknown) => {
        form.setError("cardlist", {
          message: JSON.stringify(err, null, 0),
        });
      });
  }

  return (
    <div className="p-6 flex flex-col gap-8">
      <div>
        <h2 className="scroll-m-20 text-xl font-semibold tracking-tight">
          Edit buylist
        </h2>
      </div>

      <div>
        <EditStringForm
          label="Buylist name"
          defaultValue={buylist.name}
          onSubmit={(name) => {
            UpdateBuylistName(buylist.id, name)
              .then(() => {
                void router.invalidate();
              })
              .catch((err) => {
                toast({ title: "Error", description: String(err) });
              });
          }}
        />
      </div>

      <div>
        <ButtonWithConfirm
          onClick={() => {
            DeleteBuylist(buylist.id)
              .then(() => {
                toast({
                  description: `Deleted buylist '${buylist.name}'`,
                  variant: "destructive",
                });
                void navigate({ to: "/buylists" });
              })
              .catch((err) => {
                toast({ title: "Error", description: String(err) });
              });
          }}
        ></ButtonWithConfirm>
      </div>

      <div className="flex flex-col gap-2">
        <h4 className="scroll-m-20 text-xl font-semibold tracking-tight">
          Add Cards
        </h4>
        <Form {...form}>
          <form
            autoComplete="off"
            onSubmit={(e) => {
              e.preventDefault();
              void form.handleSubmit(onSubmit)();
            }}
          >
            <FormField
              control={form.control}
              name="cardlist"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Card List</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="32 Forest"
                      className="resize-none"
                      rows={20}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Paste in the deck list from your deck builder tool
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="mt-2">
              Submit
            </Button>
          </form>
        </Form>
      </div>
    </div>
  );
}
