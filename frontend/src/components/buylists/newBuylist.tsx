import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

import { useState } from "react";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useNewBuylist } from "@/hooks/buylists";
import { BadgePlus } from "lucide-react";

const schema = z.object({
  name: z
    .string()
    .min(2, {
      message: "Name must be at least 2 characters",
    })
    .max(50, {
      message: "Name must be less than 50 characters",
    }),
});
type formSchemaType = z.infer<typeof schema>;

export const NewBuylistForm = () => {
  const form = useForm<formSchemaType>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
    },
  });

  const newBuylist = useNewBuylist();
  const onSubmit = async (values: formSchemaType) => {
    try {
      await newBuylist(values.name);
    } catch (err) {
      form.setError(
        "name",
        { message: JSON.stringify(err, null) },
        { shouldFocus: true },
      );
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          void form.handleSubmit(onSubmit)();
        }}
        autoComplete="off"
        className="space-y-8"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Buylist name" {...field} />
              </FormControl>
              <FormDescription>
                Enter a descriptive name for your buylist
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
};

export const NewBuylistDialog = () => {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          className="h-9"
          onClick={() => {
            setOpen(true);
          }}
        >
          <BadgePlus /> New
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New buylist</DialogTitle>
          <DialogDescription>Create a new buylist</DialogDescription>
        </DialogHeader>

        <NewBuylistForm />
      </DialogContent>
    </Dialog>
  );
};
