using System.Linq;

namespace Jargon.Impl
{
    internal static class KnownTokensLookup
    {
        private static readonly Token[] KnownTokens =
            {
                new Token(" ", false, true, false),
                new Token("\r", true, false, false),
                new Token("\n", true, false, false),
                new Token("\t", true, false, false),
                new Token(".", true, false, false),
                new Token(",", true, false, false)
            };

        private static readonly int Offset;
        private static readonly Token[] TokenOffsetLookup;


        static KnownTokensLookup()
        {
            var map = KnownTokens.ToDictionary(t => new Rune(t.String.Single()), t => t);

            var smallestRune = map.Min(m => m.Key.Character);
            var largestRune = map.Max(m => m.Key.Character);

            Offset = -smallestRune;
            var size = largestRune - smallestRune + 1;

            TokenOffsetLookup = new Token[size];
            foreach(var kv in map)
            {
                var pos = kv.Key.Character + Offset;
                TokenOffsetLookup[pos] = kv.Value;
            }
        }

        public static bool TryGetValue(Rune r, out Token tok)
        {
            if (r.IsMultiCharacter)
            {
                tok = default;
                return false;
            }

            var ix = r.Character + Offset;

            if (ix < 0 || ix > TokenOffsetLookup.Length)
            {
                tok = default;
                return false;
            }

            tok = TokenOffsetLookup[ix];
            return tok.String != null;
        }
    }
}
