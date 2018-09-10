using System;
using System.Collections.Generic;

namespace Jargon.Impl
{
    internal sealed class LemmaTokens<TTokenProvider> : ITokens
        where TTokenProvider: ITokens
    {
        // C#-y style bits
        public Token Current { get; private set; }

        private TTokenProvider Incoming;
        private List<Token> Buffer;
        private readonly Lemmatizer Lem;

        public LemmaTokens(Lemmatizer lem, TTokenProvider tokens)
        {
            if(tokens == null) throw new ArgumentNullException(nameof(tokens));

            Incoming = tokens;
            Lem = lem;
            Buffer = new List<Token>();
        }

        public void Dispose()
        {
            Incoming?.Dispose();
            Incoming = default(TTokenProvider);
        }

        // going closer to the Go code

        public Token? Next()
        {
            var t = this;

            while (true)
            {
                t.Fill(1);

                if(t.Buffer.Count == 0)
                {
                    return null;
                }

                var tok = t.Buffer[0];

                if(tok.IsPunct || tok.IsSpace)
                {
                    t.Drop(1);
                    return tok;
                }

                return t.NGrams();
            }
        }

        private bool Fill(int count)
        {
            var t = this;

            while (count >= t.Buffer.Count)
            {
                var token = t.Incoming.Next();
                if(token == null)
                {
                    return false;
                }

                t.Buffer.Add(token.Value);
            }

            return true;
        }

        private void Drop(int n)
        {
            var t = this;

            var toDrop = Math.Min(n, t.Buffer.Count);
            t.Buffer.RemoveRange(0, toDrop);
        }

        private Token? NGrams()
        {
            var t = this;

            for (var take = t.Lem.MaxGramLength; take > 0; take--)
            {
                var (run, consumed, ok) = t.WordRun(take);

                if (!ok) continue;

                var (canonical, found) = t.Lem.Lookup(run);

                if(found)
                {
                    var lemma = new Token(canonical, false, false, true);
                    t.Drop(consumed);
                    return lemma;
                }

                if(take == 1)
                {
                    var original = t.Buffer[0];
                    t.Drop(1);
                    return original;
                }
            }

            throw new Exception("Did not find token, this should never happen.");
        }

        private List<string> WordRun_Taken;
        internal (string[] Taken, int Count, bool Ok) WordRun(int take)
        {
            var t = this;

            var taken = (WordRun_Taken ?? (WordRun_Taken = new List<string>()));
            taken.Clear();
            var count = 0;

            while(taken.Count < take)
            {
                var ok = t.Fill(count);
                if(!ok)
                {
                    return (null, 0, false);
                }

                var token = t.Buffer[count];

                if(token.IsPunct)
                {
                    return (null, 0, false);
                }

                if (token.IsSpace)
                {
                    count++;
                    continue;
                }

                // default
                taken.Add(token.String);
                count++;
            }

            return (taken.ToArray(), count, true);
        }
    }
}
