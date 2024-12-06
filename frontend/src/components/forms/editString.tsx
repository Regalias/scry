import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Edit } from "lucide-react";
import { useState } from "react";

export interface EditStringFormProps {
  label: string;
  defaultValue: string;
  buttonText?: string;
  onSubmit: (newValue: string) => void;
}

export const EditStringForm = ({
  onSubmit,
  defaultValue,
  label,
  buttonText = "Edit",
}: EditStringFormProps) => {
  const [input, setInput] = useState(defaultValue);
  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="label">{label}</Label>
      <div className="flex w-full max-w-sm items-center space-x-2">
        <Input
          type="name"
          value={input}
          onChange={(e) => {
            setInput(e.target.value);
          }}
        />
        <Button
          type="submit"
          disabled={input === defaultValue}
          onClick={() => {
            onSubmit(input);
          }}
          variant="secondary"
        >
          <Edit />
          {buttonText}
        </Button>
      </div>
    </div>
  );
};
