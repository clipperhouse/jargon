using System;
using System.Collections;
using System.Collections.Generic;
using System.Text;

namespace Jargon
{
    public sealed class LemmaTokens : Tokens
    {
        // C#-y style bits
        public Token Current { get; private set; }

        object IEnumerator.Current => Current;

        private Tokens Incoming;
        private List<Token> Buffer;
        private readonly Lemmatizer Lem;

        public LemmaTokens(Lemmatizer lem, Tokens tokens)
        {
            if (tokens == null) throw new ArgumentNullException(nameof(tokens));

            Lem = lem;
            Incoming = tokens;
            Buffer = new List<Token>();
        }

        public void Dispose()
        {
            Incoming?.Dispose();
            Incoming = null;
        }

        public bool MoveNext()
        {
            var ret = Next();
            if(ret != null)
            {
                Current = ret.Value;
                return true;
            }

            return false;
        }

        public void Reset()
        {
            throw new NotImplementedException(nameof(Reset));
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

                if(tok.Punct || tok.Space)
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

                var gram = Join(run);
                var (canonical, found) = t.Lem.Lookup(gram);

                if(found)
                {
                    var lemma = new Token(canonical, false, false, true);
                    t.Drop(consumed);
                    return lemma;
                }

                if(take == 1)
                {
                    t.Drop(1);
                    return run[0];
                }
            }

            throw new Exception("Did not find token, this should never happen.");
        }

        internal (List<Token> Taken, int Count, bool Ok) WordRun(int take)
        {
            var t = this;

            var taken = new List<Token>();
            var count = 0;

            while(taken.Count < take)
            {
                var ok = t.Fill(count);
                if(!ok)
                {
                    return (null, 0, false);
                }

                var token = t.Buffer[count];

                if(token.Punct)
                {
                    return (null, 0, false);
                }

                if (token.Space)
                {
                    count++;
                    continue;
                }

                // default
                taken.Add(token);
                count++;
            }

            return (taken, count, true);
        }

        private string Join(List<Token> tokens)
        {
            var ret = new StringBuilder();
            foreach(var t in tokens)
            {
                ret.Append(t.Value);
            }

            return ret.ToString();
        }
    }
}
