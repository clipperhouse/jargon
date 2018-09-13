using System;

namespace Jargon.Impl
{
    internal sealed class LemmaTokens<TTokenProvider> : ITokens
        where TTokenProvider: ITokens
    {
        // C#-y style bits
        public Token Current { get; private set; }

        private TTokenProvider Incoming;
        private TokenBuffer Buffer;
        private readonly Lemmatizer Lem;
        internal WordRunBuffer WordRunBuffer;

        public LemmaTokens(Lemmatizer lem, TTokenProvider tokens)
        {
            if(tokens == null) throw new ArgumentNullException(nameof(tokens));

            Incoming = tokens;
            Lem = lem;
            Buffer = new TokenBuffer(4);
            WordRunBuffer = new WordRunBuffer(lem.MaxGramLength);
        }

        public void Dispose()
        {
            Incoming?.Dispose();
            Incoming = default;
        }

        // going closer to the Go code

        public Token? Next()
        {
            while (true)
            {
                Fill(1);

                if(Buffer.Count == 0)
                {
                    return null;
                }

                var tok = Buffer.FirstElement();

                if(tok.IsPunct || tok.IsSpace)
                {
                    Drop(1);
                    return tok;
                }

                return NGrams();
            }
        }

        private bool Fill(int count)
        {
            while (count >= Buffer.Count)
            {
                var token = Incoming.Next();
                if(token == null)
                {
                    return false;
                }

                Buffer.Add(token.Value);
            }

            return true;
        }

        private void Drop(int n)
        {
            var toDrop = Math.Min(n, Buffer.Count);
            Buffer.RemoveFromFront(toDrop);
        }

        private Token? NGrams()
        {
            for (var take = Lem.MaxGramLength; take > 0; take--)
            {
                var ok = WordRun(take, out var consumed);

                var run = WordRunBuffer.AsArray;
                var runLen = WordRunBuffer.Count;

                if (!ok) continue;

                var (canonical, found) = Lem.Lookup(run, runLen);

                if(found)
                {
                    var lemma = Token.NewLemma(canonical);
                    Drop(consumed);
                    return lemma;
                }

                if(take == 1)
                {
                    var original = Buffer.FirstElement();
                    Drop(1);
                    return original;
                }
            }

            throw new Exception("Did not find token, this should never happen.");
        }

        internal bool WordRun(int take, out int count)
        {
            WordRunBuffer.Clear();
            count = 0;

            while(WordRunBuffer.Count < take)
            {
                var ok = Fill(count);
                if(!ok)
                {
                    return false;
                }

                var token = Buffer.ElementAt(count);

                if(token.IsPunct)
                {
                    return false;
                }

                if (token.IsSpace)
                {
                    count++;
                    continue;
                }

                // default
                WordRunBuffer.Add(token.String);
                count++;
            }

            return true;
        }
    }
}
