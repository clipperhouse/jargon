using System;

namespace Jargon.Impl
{
    // basically Go's iterable interface
    internal interface ITokens: IDisposable
    {
        Token? Next();
    }
}
