import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/shadcn/ui/card";
import { cn } from "@/lib/shadcn-utils";
import { Button } from "@/components/shadcn/ui/button";
import { Input } from "@/components/shadcn/ui/input";
import { Label } from "@/components/shadcn/ui/label";
import Loader from "@/components/loader";
import { APIErrorReason, ensureAPIError } from "@/services/errors";
import { useAPIContext } from "@/contexts/api";

const ORGANIZATION_ALREADY_EXISTS = "Organization with such slug already exists."

export default function CreateOrganization({className, ...props}: React.ComponentPropsWithoutRef<"div">) {
  const navigate = useNavigate();
  const {useCreateOrganization} = useAPIContext();

  const [name, setName] = useState<string>("");
  const [slug, setSlug] = useState<string>("");

  const [nameError, setNameError] = useState<string | null>(null);
  const [slugError, setSlugError] = useState<string | null>(null);
  const [generalError, setGeneralError] = useState<string | null>(null);

  const isLoading = useCreateOrganization.isPending;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setNameError(null);
    setSlugError(null);
    setGeneralError(null);

    if (!name.trim() || !slug.trim()) {
        return;
    }

    try {
        await useCreateOrganization.mutateAsync({name: name.trim(), slug: slug.trim()});
        await navigate("/");
    } catch (e: unknown) {
        const error = ensureAPIError(e);
        if (error.reason === APIErrorReason.BAD_REQUEST) {
            setNameError(error.friendlyMessage);
            setSlugError(error.friendlyMessage);
        } else if (error.reason === APIErrorReason.ALREADY_EXISTS) {
            setSlugError(ORGANIZATION_ALREADY_EXISTS);
        } else {
            setGeneralError(error.friendlyMessage);
        }
    }
  };

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
            <div className={cn("flex flex-col gap-6", className)} {...props}>
                <Card>
                    <CardHeader>
                        <CardTitle className="text-2xl text-center">Create Organization</CardTitle>
                        <CardDescription className="text-center">
                            Enter organization details below to create an organization
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit}>
                            <div className="flex flex-col gap-6">
                                <div className="grid gap-2">
                                    <Label htmlFor="name">Name</Label>
                                    <Input
                                        id="name"
                                        type="text"
                                        required
                                        onChange={(e) => setName(e.target.value)} 
                                    />
                                    {nameError && <p className="text-red-500 text-sm text-left">{nameError}</p>}
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="slug">Slug</Label>
                                    <Input
                                        id="slug"
                                        type="text"
                                        required
                                        onChange={(e) => setSlug(e.target.value)} 
                                    />
                                    {slugError && <p className="text-red-500 text-sm text-left">{slugError}</p>}
                                </div>
                                {generalError && <p className="text-red-500 text-sm text-center">{generalError}</p>}
                                <Button type="submit" className="w-full">
                                    {isLoading ? (<><Loader className="mr-2 h-4 w-4 animate-spin" /></>) : ("Create")}
                                </Button>
                            </div>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </div>
    </div>
  );
};