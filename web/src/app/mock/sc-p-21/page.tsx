import Image from "next/image";
import { Math } from "@/components/math/Math";
import PageHeader from "../_components/PageHeader";
import ProblemTitle from "../_components/ProblemTitle";
import ProblemSectionHeader from "../_components/ProblemSectionHeader";
import ProblemParagraph from "../_components/ProblemParagraph";
import SpeechBubble from "../_components/SpeechBubble";
import WaterPump from "../_components/icons/WaterPump";
import WaterDrop from "../_components/icons/WaterDrop";
import DeleteIcon from "../_components/icons/DeleteIcon";

export default function ScreenP11() {
  return (
    <>
      <PageHeader />
      <main className="px-8 pb-8">
        <ProblemTitle>ポンプ Lv1</ProblemTitle>

        <section className="py-4">
          <ProblemSectionHeader>問題文</ProblemSectionHeader>
          <ProblemParagraph>
            M君はとあるシミュレーションゲームをプレイしています。
          </ProblemParagraph>
          <ProblemParagraph>
            そのゲームには上水道の概念があり、ポンプと水場を組み合わせることで上水道を供給できます。
          </ProblemParagraph>
          <ProblemParagraph>
            M君は縦 <Math tex="H" /> マス×横 <Math tex="W" />{" "}
            マスのグリッドからなる長方形状の土地を用意し、それぞれのマスにポンプか水場を置くことで、最大効率のポンプ場を作りたいと考えています。
          </ProblemParagraph>
          <ProblemParagraph>
            それぞれのポンプが汲み上げる水の量は以下のルールにより決まります。
          </ProblemParagraph>
          <ul className="[&>li]:py-1 [&>li]:list-disc [&>li]:ml-8">
            <li>
              周りのポンプや水場の配置によらず、基本 <Math tex="X" /> ガロン
            </li>
            <li>
              ポンプに隣接する（斜めを含む）8マスに水場があれば、追加で1マスあたり{" "}
              <Math tex="Y" /> ガロン
              <ul className="[&>li]:py-1 [&>li]:list-[circle] [&>li]:ml-8">
                <li>
                  この追加分は1つの水場の周りにポンプが何個あろうとも、それぞれのポンプに対して発生する
                </li>
              </ul>
            </li>
          </ul>
          <ProblemParagraph>
            汲み上げられる水の総量が最大になるようなポンプと水場の配置を求めてください。
          </ProblemParagraph>
        </section>

        <section className="py-4">
          <ProblemSectionHeader>制約</ProblemSectionHeader>
          <ul className="[&>li]:py-1 [&>li]:list-disc [&>li]:ml-8">
            <li>
              <Math tex="H=3, \ W=3" />
            </li>
            <li>
              <Math tex="X=10000, \ Y=1000" />
            </li>
          </ul>
        </section>

        <section className="flex flex-col gap-y-4 pt-4">
          <div className="w-full flex items-center">
            <ProblemSectionHeader className="w-full">
              解答欄
            </ProblemSectionHeader>
            <button className="btn btn-secondary">リセット</button>
          </div>
          <div className="w-full flex rounded-lg bg-base-200 p-4">
            <div className="flex-1 flex justify-center">
              <div className="flex-none grid grid-cols-3 place-items-center gap-[1px] p-[1px] bg-base-content [&>button]:bg-base-100">
                <button className="p-4 cursor-pointer">
                  <WaterPump className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <WaterPump className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <WaterPump className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <WaterPump className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <WaterDrop className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <WaterPump className="h-8" />
                </button>
                <button className="p-4 cursor-pointer">
                  <div className="h-8 w-8"></div>
                </button>
                <button className="p-4 cursor-pointer">
                  <div className="h-8 w-8"></div>
                </button>
                <button className="p-4 cursor-pointer">
                  <div className="h-8 w-8"></div>
                </button>
              </div>
            </div>
            <div className="flex-1 flex flex-col justify-between">
              <div>
                <div>パレット:</div>
                <div className="join">
                  <label className="join-item btn btn-lg has-[:checked]:btn-primary">
                    <input
                      type="radio"
                      name="palette"
                      defaultChecked
                      className="sr-only"
                    />
                    <WaterPump className="h-8 stroke-base-content" />
                  </label>
                  <label className="join-item btn btn-lg has-[:checked]:btn-primary">
                    <input type="radio" name="palette" className="sr-only" />
                    <WaterDrop className="h-8 fill-base-content" />
                  </label>
                  <label className="join-item btn btn-lg has-[:checked]:btn-primary">
                    <input type="radio" name="palette" className="sr-only" />
                    <DeleteIcon className="h-8 fill-base-content" />
                  </label>
                </div>
              </div>
              <div>
                <div className="text-xl mb-4">合計: 80000ガロン</div>
                <div className="inline-flex gap-x-4 items-center">
                  <button className="btn btn-lg btn-primary">解答する</button>
                  <span className="text-accent text-lg hidden">正解！</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>
    </>
  );
}
