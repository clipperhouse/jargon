using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace Jargon.Impl
{
    internal sealed class TextTokens : ITokens
    {
        // C#-y bits
        
        private TextReader Incoming;
        private Rune? PendingRune;
        private CharBuffer Buffer;

        public TextTokens(TextReader reader)
        {
            Incoming = reader ?? throw new ArgumentNullException(nameof(reader));
            Buffer = CharBuffer.New();
        }

        public void Dispose()
        {
            Incoming?.Dispose();
            Incoming = null;
        }

        // doing the Go stuff that C# doesn't have a great analog for

        private bool TryReadRune(out Rune r) // returns null if we've hit the end of the reader
        {
            if (PendingRune.HasValue)
            {
                var ret = PendingRune.Value;
                PendingRune = null;
                r = ret;
                return true;
            }

            var reader = Incoming;
            var resHigh = reader.Read();
            if (resHigh == -1)
            {
                // eof
                r = default;
                return false;
            }

            var cHigh = (char)resHigh;
            if (!char.IsHighSurrogate(cHigh))
            {
                // simple character
                r = new Rune(cHigh);
                return true;
            }

            var resLow = reader.Read();
            if (resLow == -1)
            {
                // eof, with an invalid char before it
                r = default;
                return false;
            }

            var cLow = (char)resLow;
            if (!char.IsLowSurrogate(cLow))
            {
                // not paired, not end of file, so skip it and read the next one
                return TryReadRune(out r);
            }

            // paired surrogates character
            r = new Rune(cHigh, cLow);
            return true;
        }

        private void UnreadRune(Rune r)
        {
            if (PendingRune != null)
            {
                throw new InvalidOperationException($"Called {nameof(UnreadRune)} when pending rune already set");
            }

            PendingRune = r;
        }

        // going closer to the Go code

        public Token? Next()
        {
            if (Buffer.Length > 0)
            {
                return Token();
            }

            while (true)
            {
                if (!TryReadRune(out var r))
                {
                    // end of the input
                    return Token();
                }

                if (r.IsSpace())
                {
                    Accept(r);
                    return Token();
                }

                if (r.IsPunct())
                {
                    Accept(r);
                    var isLeadingPunct = r.MightBeLeadingPunct() && !PeekTerminator();
                    if (isLeadingPunct)
                    {
                        return ReadWord();
                    }

                    return Token();
                }

                Accept(r);
                return ReadWord();
            }
        }

        private Token? ReadWord()
        {
            while (true)
            {
                if (!TryReadRune(out var r))
                {
                    return Token();
                }

                if (r.MightBeMidPunct())
                {
                    var followedByTerminator = PeekTerminator();
                    if (followedByTerminator)
                    {
                        var tok = Token();
                        Accept(r);

                        return tok;
                    }

                    Accept(r);
                    continue;
                }

                if (r.IsPunct() || r.IsSpace())
                {
                    var tok = Token();

                    Accept(r);

                    return tok;
                }

                Accept(r);
            }
        }

        private Token? Token()
        {
            if (Buffer.Length == 0)
            {
                return null;
            }

            if (Buffer.IsSingleRune(out var r))
            {
                Buffer.Reset();

                if (KnownTokensLookup.TryGetValue(r, out var known))
                {
                    return known;
                }

                return new Token(r.AsString, r.IsPunct(), r.IsSpace(), false);
            }

            var b = Buffer.ToStringReset();
            return new Token(b, false, false, false);
        }

        private void Accept(Rune r)
        {
            Buffer.Append(r);
        }

        private bool PeekTerminator()
        {
            if (!TryReadRune(out var r))
            {
                return true;
            }

            UnreadRune(r);

            return r.IsPunct() || r.IsSpace();
        }
    }
}