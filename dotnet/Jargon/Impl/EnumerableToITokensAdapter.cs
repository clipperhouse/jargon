using System;
using System.Collections.Generic;

namespace Jargon.Impl
{
    internal sealed class EnumerableToITokensAdapter : ITokens
    {
        private IEnumerator<Token> Inner;

        public Token Current => Inner.Current;

        internal EnumerableToITokensAdapter(IEnumerator<Token> tok)
        {
            Inner = tok;
        }

        public Token? Next()
        {
            var res = Inner.MoveNext();
            if (!res) return null;

            return Inner.Current;
        }

        public void Dispose()
        {
            Inner?.Dispose();
            Inner = null;
        }
    }
}
