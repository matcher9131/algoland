import Link from "next/link";
import MenuIcon from "../../../components/icons/MenuIcon";

export default function PageHeader() {
    return (
        <header className="w-full p-8 flex justify-between">
            <h1 className="text-xl">algoland</h1>
            <button className="btn btn-circle btn-ghost w-12 h-12" popoverTarget="popover_menu" style={{ anchorName: "--anchor-menu" }}>
                <MenuIcon />
            </button>
            <ul
                className="dropdown dropdown-end menu w-52 rounded-box bg-base-100 shadow-sm"
                popover="auto"
                id="popover_menu"
                style={{ positionAnchor: "--anchor-menu"}}
            >
                <li>
                    <button>ログアウト</button>
                </li>
            </ul>
        </header>
    );
}
