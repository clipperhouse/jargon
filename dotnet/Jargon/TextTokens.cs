using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace Jargon
{
    public sealed class TextTokens : Tokens
    {
        // C#-y bits
        private static readonly Dictionary<Rune, Token> KnownTokens =
            new Dictionary<Rune, Token>
            {
                [new Rune(' ')] = new Token(" ", false, true, false),
                [new Rune('\r')] = new Token("\r", true, false, false),
                [new Rune('\n')] = new Token("\n", true, false, false),
                [new Rune('\t')] = new Token("\t", true, false, false),
                [new Rune('.')] = new Token(".", true, false, false),
                [new Rune(',')] = new Token(",", true, false, false)
            };

        private TextReader Incoming;
        private Rune? PendingRune;
        private StringBuilder Buffer;

        public Token Current { get; private set; }

        object IEnumerator.Current => Current;

        public TextTokens(TextReader reader)
        {
            Incoming = reader ?? throw new ArgumentNullException(nameof(reader));
            Buffer = new StringBuilder();
        }

        public bool MoveNext()
        {
            var next = Next();
            if (next != null)
            {
                Current = next.Value;
                return true;
            }

            return false;
        }

        public void Dispose()
        {
            Incoming?.Dispose();
            Incoming = null;
        }

        public void Reset()
        {
            throw new NotImplementedException(nameof(Reset));
        }
        
        // doing the Go stuff that C# doesn't have a great analog for

        private struct Rune
        {
            public bool IsMultiCharacter => _HighSurrogate.HasValue && _LowSurrogate.HasValue;

            private string _AsString;
            public string AsString
            {
                get
                {
                    if (_AsString != null)
                    {
                        return _AsString;
                    }

                    if (IsMultiCharacter)
                    {
                        _AsString = char.ConvertFromUtf32(char.ConvertToUtf32(HighSurrogate, LowSurrogate));
                    }
                    else
                    {
                        _AsString = Character.ToString();
                    }
                    return _AsString;
                }
            }

            public char Character
            {
                get
                {
                    if (IsMultiCharacter) throw new InvalidOperationException($"Cannot access multi-char {nameof(Rune)} with {nameof(Character)}");

                    return _HighSurrogate.Value;
                }
            }

            private readonly char? _HighSurrogate;
            public char HighSurrogate
            {
                get
                {
                    if (!IsMultiCharacter) throw new InvalidOperationException($"Cannot access single-char {nameof(Rune)} with {nameof(HighSurrogate)}");

                    return _HighSurrogate.Value;
                }
            }
            private readonly char? _LowSurrogate;
            public char LowSurrogate
            {
                get
                {
                    if (!IsMultiCharacter) throw new InvalidOperationException($"Cannot access single-char {nameof(Rune)} with {nameof(LowSurrogate)}");

                    return _LowSurrogate.Value;
                }
            }

            public Rune(char c)
            {
                if (char.IsHighSurrogate(c) || char.IsLowSurrogate(c)) throw new ArgumentException($"Cannot create single-char {nameof(Rune)} with surrogate char", nameof(c));

                _HighSurrogate = c;
                _LowSurrogate = null;
                _AsString = null;
            }

            public Rune(char high, char low)
            {
                if (!char.IsHighSurrogate(high)) throw new ArgumentException($"Cannot create multi-char {nameof(Rune)} with non-high-surrogate char", nameof(high));
                if (!char.IsLowSurrogate(low)) throw new ArgumentException($"Cannot create multi-char {nameof(Rune)} with non-high-surrogate char", nameof(low));

                _HighSurrogate = high;
                _LowSurrogate = low;
                _AsString = null;
            }

            public bool IsSpace()
            {
                if (IsMultiCharacter)
                {
                    return char.IsWhiteSpace(AsString, 0);
                }

                return char.IsWhiteSpace(Character);
            }

            public bool MightBeLeadingPunct()
            {
                if (IsMultiCharacter) return false;

                // only case
                return Character == '.';
            }

            public bool MightBeMidPunct()
            {
                if (IsMultiCharacter) return false;

                switch (Character)
                {
                    case '.':
                    case '\'':
                    case '’':
                    case ':':
                    case '?':
                    case '&':
                        return true;
                }

                return false;
            }

            public bool IsPunct()
            {
                if (IsMultiCharacter)
                {
                    return char.IsPunctuation(AsString, 0); // none of the exceptions are multi-char, so elide the call
                }

                return char.IsPunctuation(Character) && !IsPunctException();
            }

            private bool IsPunctException()
            {
                if (IsMultiCharacter) return false;

                switch (Character)
                {
                    case '-':
                    case '#':
                    case '@':
                    case '*':
                    case '%':
                    case '_':
                    case '/':
                    case '\\': return true;
                }

                return false;
            }

            public static bool IsSingleRune(string str)
            {
                if (str.Length == 0) return false;
                if (str.Length == 1) return true;
                if (str.Length > 2) return false;

                return char.IsHighSurrogate(str[0]) && char.IsLowSurrogate(str[1]);
            }
        }

        private (Rune Rune, bool EndOfFile) ReadRune()
        {
            if (PendingRune.HasValue)
            {
                var ret = (PendingRune.Value, false);
                PendingRune = null;
                return ret;
            }

            var reader = Incoming;
            var resHigh = reader.Read();
            if (resHigh == -1) return (default(Rune), true);

            var cHigh = (char)resHigh;
            if (!char.IsHighSurrogate(cHigh))
            {
                return (new Rune(cHigh), false);
            }

            var resLow = reader.Read();
            if (resLow == -1) return (default(Rune), true);

            var cLow = (char)resLow;
            if (!char.IsLowSurrogate(cLow))
            {
                return (default(Rune), false);
            }

            return (new Rune(cHigh, cLow), false);
        }

        private void UnreadRune(Rune r)
        {
            if (PendingRune != null)
            {
                throw new InvalidOperationException($"Called {nameof(UnreadRune)} when pending rune already set");
            }

            PendingRune = r;
        }

        private static void WriteRune(StringBuilder b, Rune rune)
        {
            if (rune.IsMultiCharacter)
            {
                b.Append(rune.HighSurrogate);
                b.Append(rune.LowSurrogate);
            }
            else
            {
                b.Append(rune.Character);
            }
        }

        // going closer to the Go code

        public Token? Next()
        {
            var t = this;

            if (Buffer.Length > 0)
            {
                return t.Token();
            }

            while (true)
            {
                var (r, eof) = ReadRune();
                if (eof)
                {
                    // end of the input
                    return t.Token();
                }

                if (r.IsSpace())
                {
                    t.Accept(r);
                    return t.Token();
                }

                if (r.IsPunct())
                {
                    t.Accept(r);
                    var isLeadingPunct = r.MightBeLeadingPunct() && !t.PeekTerminator();
                    if (isLeadingPunct)
                    {
                        return t.ReadWord();
                    }

                    return t.Token();
                }

                t.Accept(r);
                return t.ReadWord();
            }
        }

        private Token? ReadWord()
        {
            var t = this;
            while (true)
            {
                var (r, eof) = ReadRune();
                if (eof)
                {
                    return t.Token();
                }

                if (r.MightBeMidPunct())
                {
                    var followedByTerminator = t.PeekTerminator();
                    if (followedByTerminator)
                    {
                        var tok = t.Token();
                        t.Accept(r);

                        return tok;
                    }

                    t.Accept(r);
                    continue;
                }

                if (r.IsPunct() || r.IsSpace())
                {
                    var tok = t.Token();

                    t.Accept(r);

                    return tok;
                }

                t.Accept(r);
            }
        }

        private Token? Token()
        {
            var t = this;

            var b = t.Buffer.ToString();
            if (b.Length == 0)
            {
                return null;
            }

            t.Buffer.Clear();

            if (Rune.IsSingleRune(b))
            {
                var r = b.Length == 1 ? new Rune(b[0]) : new Rune(b[0], b[1]);

                if (KnownTokens.TryGetValue(r, out var known))
                {
                    return known;
                }

                return new Token(r.AsString, r.IsPunct(), r.IsSpace(), false);
            }

            return new Token(b, false, false, false);
        }

        private void Accept(Rune r)
        {
            var t = this;
            WriteRune(t.Buffer, r);
        }

        private bool PeekTerminator()
        {
            var t = this;

            var (r, eof) = ReadRune();

            if (eof)
            {
                return true;
            }

            UnreadRune(r);

            return r.IsPunct() || r.IsSpace();
        }
    }
}