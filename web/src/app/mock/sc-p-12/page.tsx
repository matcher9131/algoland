import Image from "next/image";
import { Math } from "@/components/math/Math";
import PageHeader from "../_components/PageHeader";
import ProblemTitle from "../_components/ProblemTitle";
import ProblemSectionHeader from "../_components/ProblemSectionHeader";
import ProblemParagraph from "../_components/ProblemParagraph";
import SpeechBubble from "../_components/SpeechBubble";

export default function ScreenP12() {
  return (
    <>
      <PageHeader />
      <main className="px-8 pb-8">
        <ProblemTitle>呑兵衛 Lv2</ProblemTitle>

        <section className="py-4">
          <ProblemSectionHeader>問題文</ProblemSectionHeader>
          <ProblemParagraph className="text-warning">
            ※ 問題文は呑兵衛Lv1と同様です。制約のみが変化しています。
          </ProblemParagraph>
          <ProblemParagraph>
            大酒飲みのAさんは、アルコール度数が
            <Math tex="X" />
            %以上のお酒はおいしく感じますが、逆に
            <Math tex="X" />
            %未満のお酒は不味いと感じます。ただし、本人にその自覚はないようです。
          </ProblemParagraph>
          <ProblemParagraph>
            Aさんの部下のBさんはそのことを知っていますが、 <Math tex="X" />
            の値までは知りません。
          </ProblemParagraph>
          <ProblemParagraph>
            今日はAさんの誕生日です。上司思いのBさんは、Aさんにお酒を飲ませて気分良く酔ってもらいたいと思っています。一方で、必要以上に度数の高いお酒を飲ませてAさんの肝臓に負担をかけたくありません。
          </ProblemParagraph>
          <ProblemParagraph>
            そこで、今日こそは
            <Math tex="X" />
            の値を解明したいと思い、様々な度数のお酒を用意しました。
          </ProblemParagraph>
          <ProblemParagraph>
            部下への信頼が厚いAさんはBさんの勧めるがままにお酒を飲みます。あなたはBさんになりきって、できるだけ少ない回数で
            <Math tex="X" />
            を特定してください。
          </ProblemParagraph>
        </section>

        <section className="py-4">
          <ProblemSectionHeader>制約</ProblemSectionHeader>
          <ul className="[&>li]:py-1 [&>li]:list-disc [&>li]:ml-8">
            <li>
              <Math tex="X" />
              は実数（倍精度浮動小数点数で表せる数）
            </li>
            <li>
              <Math tex="0 \lt X \leq 100" />
            </li>
            <li>
              17回以内であり得る
              <Math tex="X" />
              の範囲が0.001(%)以下になれば正解、そうでなければ不正解とする
            </li>
            <li>入力は0より大きく100以下の10進数のみ受け付ける</li>
          </ul>
        </section>

        <section className="flex flex-col gap-y-4 pt-4">
          <div className="w-full flex items-center">
            <ProblemSectionHeader className="w-full">
              解答欄
            </ProblemSectionHeader>
            <button className="btn btn-secondary">新しく始める</button>
          </div>
          <div className="w-full p-4 flex flex-col rounded-lg bg-base-200">
            <div className="w-full py-2 flex justify-center">
              <div className="badge badge-lg badge-info">
                あなたはAさんに度数
                <span className="font-bold text-lg underline">90%</span>
                のお酒を飲ませた。
              </div>
            </div>
            <div className="w-full flex flex-row items-center">
              <div className="flex flex-col items-center">
                <Image
                  src="/figure_dance.png"
                  alt="うれしい"
                  width={70}
                  height={70}
                />
                <div className="text-xs">Aさん</div>
              </div>
              <SpeechBubble>おいしい！</SpeechBubble>
            </div>
            <div className="w-full py-2 flex justify-center">
              <div className="badge badge-lg badge-info">
                あなたはAさんに度数
                <span className="font-bold text-lg underline">10%</span>
                のお酒を飲ませた。
              </div>
            </div>
            <div className="w-full flex flex-row items-center">
              <div className="flex flex-col items-center">
                <Image
                  src="/figure_zasetsu.png"
                  alt="かなしい"
                  width={70}
                  height={70}
                />
                <div className="text-xs">Aさん</div>
              </div>
              <SpeechBubble>おいしくない……</SpeechBubble>
            </div>
            <div className="w-full py-2 flex justify-center">
              <div className="badge badge-lg badge-info">
                あなたはAさんに度数
                <span className="font-bold text-lg underline">80%</span>
                のお酒を飲ませた。
              </div>
            </div>
            <div className="w-full flex flex-row items-center">
              <div className="flex flex-col items-center">
                <Image
                  src="/figure_dance.png"
                  alt="うれしい"
                  width={70}
                  height={70}
                />
                <div className="text-xs">Aさん</div>
              </div>
              <SpeechBubble>おいしい！</SpeechBubble>
            </div>
            <div className="w-full py-2 flex justify-center">
              <div className="badge badge-lg badge-info">
                あなたはAさんに度数
                <span className="font-bold text-lg underline">20%</span>
                のお酒を飲ませた。
              </div>
            </div>
            <div className="w-full flex flex-row items-center">
              <div className="flex flex-col items-center">
                <Image
                  src="/figure_zasetsu.png"
                  alt="かなしい"
                  width={70}
                  height={70}
                />
                <div className="text-xs">Aさん</div>
              </div>
              <SpeechBubble>おいしくない……</SpeechBubble>
            </div>
            <div className="w-full py-2 flex justify-center">
              <div className="badge badge-xl badge-error">不正解……</div>
            </div>
          </div>
          <div className="w-full grid grid-cols-[2fr_1fr] gap-x-4 items-center">
            <input
              type="number"
              required
              min="0"
              max="100"
              step="any"
              placeholder="度数を入力"
              className="w-full input input-lg validator text-right [&::-webkit-inner-spin-button]:ml-2"
            />
            <button type="button" className="btn btn-primary btn-lg">
              飲ませる！
            </button>
            <p className="validator-hint">
              1以上100以下の整数を入力してください
            </p>
            <div className="text-right">残り: 3回</div>
          </div>
        </section>
      </main>
    </>
  );
}
