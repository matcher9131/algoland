import Link from "next/link";

export default function MockRoot() {
  return (
    <div className="flex flex-col">
      <Link href="mock/sc-p-11" className="text-xl underline">
        sc-p-11
      </Link>
      <Link href="mock/sc-p-12" className="text-xl underline">
        sc-p-12
      </Link>
      <Link href="mock/sc-p-21" className="text-xl underline">
        sc-p-21
      </Link>
    </div>
  );
}
