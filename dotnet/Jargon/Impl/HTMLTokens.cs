using HtmlAgilityPack;
using System;
using System.Collections.Generic;
using System.IO;

namespace Jargon.Impl
{
    internal sealed class HTMLTokens : ITokens
    {
        // C#-y style bits
        
        private TextTokens Text;
        private IEnumerator<(string Text, bool IsTag)> Html;

        public HTMLTokens(TextReader reader)
        {
            if (reader == null) throw new ArgumentNullException(nameof(reader));

            var doc = new HtmlDocument();
            doc.Load(reader);
            Html = GetHtmlTokenizedEnumerator(doc.DocumentNode);
        }

        public void Dispose()
        {
            Text?.Dispose();
            Text = null;

            Html?.Dispose();
            Html = null;
        }

        // have to fake how Go's html parser behaves

        private static IEnumerator<(string Text, bool IsTag)> GetHtmlTokenizedEnumerator(HtmlNode html)
        {
            if (html.NodeType == HtmlNodeType.Comment)
            {
                yield return (html.OuterHtml, true);
                yield break;
            }

            if (html.NodeType == HtmlNodeType.Text)
            {
                yield return (html.InnerText, false);
                yield break;
            }

            if (html.NodeType == HtmlNodeType.Document)
            {
                if (html.HasChildNodes)
                {
                    foreach (var child in html.ChildNodes)
                    {
                        using (var childE = GetHtmlTokenizedEnumerator(child))
                        {
                            while (childE.MoveNext()) yield return childE.Current;
                        }
                    }
                }

                yield break;
            }

            if(html.NodeType == HtmlNodeType.Element)
            {
                var outer = html.OuterHtml;
                var inner = html.InnerHtml;

                var endStartTag = outer.IndexOf(inner);

                if (endStartTag != -1) {
                    var startTag= outer.Substring(0, endStartTag);

                    if (!string.IsNullOrWhiteSpace(startTag))
                    {
                        yield return (startTag, true);
                    }
                }

                if (html.HasChildNodes)
                {
                    foreach (var child in html.ChildNodes)
                    {
                        using (var childE = GetHtmlTokenizedEnumerator(child))
                        {
                            while (childE.MoveNext()) yield return childE.Current;
                        }
                    }
                }

                var startEndTag = endStartTag == -1 ? inner.Length : endStartTag + inner.Length;
                var endTag = outer.Substring(startEndTag);
                if (!string.IsNullOrWhiteSpace(endTag))
                {
                    yield return (endTag, true);
                }
            }
        }

        // going closer to the Go code

        public Token? Next()
        {
            var t = this;

            var text = t.Text?.Next();

            if (text != null)
            {
                return text;
            }

            while (true)
            {
                if (!t.Html.MoveNext())
                {
                    return null;
                }

                var tok = t.Html.Current;
                if (!tok.IsTag)
                {
                    if(t.Text != null)
                    {
                        t.Text.Dispose();
                        t.Text = null;
                    }

                    var r = new StringReader(tok.Text);
                    t.Text = new TextTokens(r);
                    return t.Text.Next();
                }

                if (tok.IsTag)
                {
                    return Token.NewPunct(tok.Text);
                }
            }
        }
    }
}
