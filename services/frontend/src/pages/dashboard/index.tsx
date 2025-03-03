import { useUserContext } from "@/contexts/user";

export default function Dashboard({className, ...props}: React.ComponentPropsWithoutRef<"div">) {
  const {user, organizationMemberships} = useUserContext();

  return (
    <div>
        <p>User: {JSON.stringify(user)}</p>
        <p>Memberships: {JSON.stringify(organizationMemberships)}</p>
    </div>
  )
};