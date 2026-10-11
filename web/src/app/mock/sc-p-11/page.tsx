import { Math } from "@/components/math/Math";
import PageHeader from "../_components/PageHeader";
import ProblemTitle from "../_components/ProblemTitle";
import ProblemSectionHeader from "../_components/ProblemSectionHeader";
import ProblemParagraph from "../_components/ProblemParagraph";

export default function ScreenP11() {
    return (
        <>
            <PageHeader />
            <main>
                <ProblemTitle>呑兵衛 Lv1</ProblemTitle>

                <ProblemSectionHeader>問題文</ProblemSectionHeader>
                <ProblemParagraph>大酒飲みのAさんは、アルコール度数が<Math tex="X" />%以上のお酒はおいしく感じますが、逆に<Math tex="X" />%未満のお酒は不味いと感じます。ただし、本人にその自覚はないようです。</ProblemParagraph>
                <ProblemParagraph>Aさんの部下のBさんはそのことを知っていますが、 <Math tex="X" />の値までは知りません。</ProblemParagraph>
                <ProblemParagraph>今日はAさんの誕生日です。上司思いのBさんは、Aさんにお酒を飲ませて気分良く酔ってもらいたいと思っています。一方で、必要以上に度数の高いお酒を飲ませてAさんの肝臓に負担をかけたくありません。</ProblemParagraph>
                <ProblemParagraph>そこで、今日こそは<Math tex="X" />の値を解明したいと思い、様々な度数のお酒を用意しました。</ProblemParagraph>
                <ProblemParagraph>部下への信頼が厚いAさんはBさんの勧めるがままにお酒を飲みます。あなたはBさんになりきって、できるだけ少ない回数で<Math tex="X" />を特定してください。</ProblemParagraph>
                <ProblemParagraph>※ まぐれ当たりを防止するため、<Math tex="X" />は固定ではなく矛盾しない範囲で変動します。</ProblemParagraph>

                <ProblemSectionHeader>制約</ProblemSectionHeader>
                <ul className="[&>li]:py-1">
                    <li><Math tex="X" />は整数</li>
                    <li><Math tex="1 \leq X \leq 100" /></li>
                    <li>8回以内で特定できなければ不正解とする</li>
                    <li>入力は1以上100以下の10進整数のみ受け付ける</li>
                </ul>
            </main>
        </>
    )
}